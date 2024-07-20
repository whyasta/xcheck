package utils

import (
	"bigmind/xcheck-be/config"
	"bigmind/xcheck-be/internal/constant/response"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PagingInfo struct {
	Total int `json:"total"`
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

type MetaResponse struct {
	PagingInfo PagingInfo `json:"paging_info"`
}

// swagger:response APIResponse
type APIResponse[T any] struct {
	Code         int           `json:"code"`
	Status       string        `json:"status"`
	Message      string        `json:"message,omitempty"`
	Data         interface{}   `json:"data,omitempty"`
	Token        string        `json:"access_token,omitempty"`
	RefreshToken string        `json:"refresh_token,omitempty"`
	Meta         *MetaResponse `json:"meta,omitempty"`
	ResponseTime float64       `json:"time,omitempty"`
}

// swagger:parameters  IDParam
type IDParam struct {
	// The ID of the order
	//
	// in: path
	// required: true
	ID int64 `json:"id"`
}

func Null() interface{} {
	return nil
}

func BuildResponse[T any](code int, responseStatus response.ResponseStatus, message string, data T) APIResponse[T] {
	return BuildResponse_(code, responseStatus.GetResponseStatus(), message, data, nil)
}

func BuildResponseWithPaginate[T any](code int, responseStatus response.ResponseStatus, message string, data T, meta *MetaResponse) APIResponse[T] {
	return BuildResponse_(code, responseStatus.GetResponseStatus(), message, data, meta)
}

func BuildResponseWithToken[T any](code int, responseStatus response.ResponseStatus, token string, refreshToken string, message string, data T) APIResponse[T] {
	return BuildResponseWithToken_(code, responseStatus.GetResponseStatus(), token, refreshToken, message, data)
}

func BuildResponse_[T any](code int, status string, message string, data T, meta *MetaResponse) APIResponse[T] {
	return APIResponse[T]{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}

func BuildResponseWithToken_[T any](code int, status string, token string, refreshToken string, message string, data T) APIResponse[T] {
	return APIResponse[T]{
		Code:         code,
		Status:       status,
		Message:      message,
		Data:         data,
		Token:        token,
		RefreshToken: refreshToken,
	}
}

func PanicException_(key string, message string) {
	err := errors.New(strings.ReplaceAll(message, ":", " -"))
	err = fmt.Errorf("%s: %w", key, err)
	if err != nil {
		panic(err)
	}
}

func PanicException(responseKey response.ResponseStatus, message string) {
	PanicException_(responseKey.GetResponseStatus(), message)
}

func ResponseHandler(c *gin.Context) {
	if err := recover(); err != nil {
		str := fmt.Sprint(err)
		strArr := strings.Split(str, ":")

		key := strArr[0]
		msg := strings.Trim(strArr[1], " ")

		switch key {
		case
			response.DataNotFound.GetResponseStatus():
			c.JSON(http.StatusNotFound, BuildResponse_(http.StatusNotFound, key, msg, Null(), nil))
			c.Abort()
		case
			response.InvalidRequest.GetResponseStatus():
			c.JSON(http.StatusBadRequest, BuildResponse_(http.StatusBadRequest, key, msg, Null(), nil))
			c.Abort()
		case
			response.Unauthorized.GetResponseStatus():
			c.JSON(http.StatusUnauthorized, BuildResponse_(http.StatusUnauthorized, key, msg, Null(), nil))
			c.Abort()
		default:
			c.JSON(http.StatusInternalServerError, BuildResponse_(http.StatusInternalServerError, key, msg, Null(), nil))
			c.Abort()
		}
	}
}

func WriterHandler(c *gin.Context) {
	w := NewResponseWriter(c)
	c.Writer = w
	defer w.Done(c)

	t := time.Now()

	c.Next()

	latency := time.Since(t).Seconds()

	originalBody := w.body
	w.body = &bytes.Buffer{}
	var body map[string]interface{}
	json.Unmarshal(originalBody.Bytes(), &body)
	body["code"] = c.Writer.Status()
	body["processing_time"] = latency
	newBody, _ := json.Marshal(body)

	// log.Printf("Response :" + string(newBody) + "\n")

	fmt := fmt.Sprintf("%d %s %s",
		c.Writer.Status(),
		c.Request.Method,
		c.Request.RequestURI,
	)

	logData := map[string]interface{}{
		"request": map[string]interface{}{
			"headers": c.Request.Header,
			"body":    c.Request.Body,
		},
		"response": map[string]interface{}{
			"headers": c.Writer.Header(),
			"body":    json.RawMessage(newBody),
		},
	}
	config.Logger.
		WithOptions(zap.Fields(zap.Any("context", logData))).
		Info(fmt)

		// Logger.Info("fmt", zap.Fields(zap.String("context", string(newBody))))
		// Logger.Info(fmt)

	w.Write(newBody)
}
