package controllers

import (
	"bigmind/xcheck-be/checks"
	"bigmind/xcheck-be/config"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type CheckStatus struct {
	Name string `json:"name"`
	Pass bool   `json:"pass"`
}

type HealthController struct{}

func (h HealthController) Init(c *gin.Context) {
	// config.GetEnqueuer().Enqueue("test", work.Q{
	// 	"email_address": "qjDpS@example.com",
	// })
	c.JSON(http.StatusOK, gin.H{"message": "pong from " + config.GetConfig().GetString("APP_ENV")})
}

func (h HealthController) Timeout(c *gin.Context) {
	timer := time.NewTimer(34 * time.Second)

	fmt.Println("Timer started, waiting...")

	// Wait until the timer fires
	<-timer.C
	c.JSON(http.StatusOK, gin.H{"message": "timeout from " + config.GetConfig().GetString("APP_ENV")})
}

func (h HealthController) Status(checks []checks.Check, failureNotification *checks.FailureNotification) gin.HandlerFunc {
	var lock sync.Mutex
	var failureInARow uint32

	fn := func(c *gin.Context) {
		var eg errgroup.Group

		statuses := make([]CheckStatus, len(checks))
		httpStatus := http.StatusOK
		for idx, check := range checks {
			captureCheck := check
			captureIdx := idx
			eg.Go(func() error {
				pass := captureCheck.Pass()
				statuses[captureIdx] = CheckStatus{
					Name: captureCheck.Name(),
					Pass: pass,
				}

				if !pass {
					return errors.New("healthcheck failed")
				}
				return nil
			})
		}

		lock.Lock()
		if err := eg.Wait(); err != nil {
			httpStatus = http.StatusInternalServerError
			failureInARow += 1

			if failureInARow >= 1 &&
				failureNotification.Chan != nil {
				failureNotification.Chan <- err
			}
		} else {
			if failureInARow != 0 && failureNotification.Chan != nil {
				failureInARow = 0
				failureNotification.Chan <- nil
			}
		}
		lock.Unlock()

		c.JSON(httpStatus, gin.H{"message": "OK", "data": statuses})
	}

	return gin.HandlerFunc(fn)

	// c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
