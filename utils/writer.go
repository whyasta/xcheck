//lint:file-ignore ST1016 empty block comments are valid in this file
package utils

import (
	"bytes"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/gin-gonic/gin"

	"github.com/vearne/gin-timeout/buffpool"
)

type ResponseWriter struct {
	gin.ResponseWriter
	// header
	h http.Header
	// body
	body *bytes.Buffer

	code        int
	mu          sync.Mutex
	timedOut    atomic.Bool
	wroteHeader atomic.Bool
	size        int
}

func NewResponseWriter(c *gin.Context) *ResponseWriter {
	buffer := buffpool.GetBuff()
	writer := &ResponseWriter{
		body:           buffer,
		ResponseWriter: c.Writer,
		h:              make(http.Header),
	}
	return writer
}

func (tw *ResponseWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut.Load() {
		return 0, nil
	}
	tw.size += len(b)
	return tw.body.Write(b)
}

func (tw *ResponseWriter) WriteHeader(code int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut.Load() {
		return
	}
	tw.writeHeader(code)
}

func (tw *ResponseWriter) writeHeader(code int) {
	tw.wroteHeader.Store(true)
	tw.code = code
}

func (tw *ResponseWriter) WriteHeaderNow() {}

func (tw *ResponseWriter) Header() http.Header {
	return tw.h
}

func (tw *ResponseWriter) Size() int {
	return tw.size
}

func (tw *ResponseWriter) Status() int {
	if tw.code == 0 || !tw.wroteHeader.Load() {
		return tw.ResponseWriter.Status()
	}
	return tw.code
}

func (w *ResponseWriter) Done(c *gin.Context) {
	dst := w.ResponseWriter.Header()
	for k, vv := range w.Header() {
		dst[k] = vv
	}

	if !w.wroteHeader.Load() {
		w.code = http.StatusOK
	}

	w.ResponseWriter.WriteHeader(w.code)
	_, err := w.ResponseWriter.Write(w.body.Bytes())
	if err != nil {
		panic(err)
	}
	buffpool.PutBuff(w.body)
}
