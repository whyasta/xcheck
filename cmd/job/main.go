package main

import (
	"bigmind/xcheck-be/config"
	"bigmind/xcheck-be/internal/processors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gocraft/work"
)

type Context struct {
	Id   int64
	Data map[string]interface{}
}

// var redisPool = &redis.Pool{
// 	MaxActive: 5,
// 	MaxIdle:   5,
// 	Wait:      true,
// 	Dial: func() (redis.Conn, error) {
// 		redisHost := config.GetConfig().GetString("REDIS_HOST")
// 		redisPort := config.GetConfig().GetString("REDIS_PORT")
// 		return redis.Dial("tcp", fmt.Sprintf("%s:%s", redisHost, redisPort))
// 	},
// }

func main() {
	fmt.Println("Running job queue")
	config.Init("production")

	pool := work.NewWorkerPool(Context{}, 10, "xcheck", config.NewRedis())
	pool.Middleware((*Context).Log)
	pool.Start()

	pool.Job("test", processors.TestJob)
	pool.Job("import_barcode", processors.ImportBarcodeJob)

	// waiting exit signal：
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	fmt.Println("Stop the pool")

	// Stop the pool
	pool.Stop()
}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	log.Println("Starting job: ", job.Name)
	return next()
}
