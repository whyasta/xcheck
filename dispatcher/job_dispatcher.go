package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type Job struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

func dispatchJobs(rdb *redis.Client, queueName string, numJobs int) {
	for j := 1; j <= numJobs; j++ {
		jobID := fmt.Sprintf("job-%d", j)
		err := rdb.RPush(ctx, queueName, jobID).Err()
		if err != nil {
			log.Printf("Error pushing job to queue: %v", err)
			continue
		}
		fmt.Printf("Dispatched job: %s\n", jobID)
	}
}

func dispatchJob(rdb *redis.Client, queueName string, job Job) {
	// jobID := fmt.Sprintf("job-%d", jobID)

	b, err := json.Marshal(job)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = rdb.RPush(ctx, queueName, string(b)).Err()
	if err != nil {
		log.Printf("Error pushing job to queue: %v", err)
		return
	}
	fmt.Printf("Dispatched job: %v\n", job.Type)
}

func dispatchJob1(rdb *redis.Client, queueName string, jobID int) {
	// jobID := fmt.Sprintf("job-%d", jobID)
	err := rdb.RPush(ctx, queueName, jobID).Err()
	if err != nil {
		log.Printf("Error pushing job to queue: %v", err)
		return
	}
	fmt.Printf("Dispatched job: %v\n", jobID)
}

func main() {
	// Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Change if Redis is on another host
		DB:   0,                // Use default DB
	})

	// Redis queue name
	queueName := "workerQueue"
	// numJobs := 2

	dispatchJob(rdb, queueName, Job{
		Type: "disbursement",
		Data: map[string]interface{}{"key": "value"},
	})

	//dispatchJob1(rdb, queueName, 84055)

	// dispatchJobs(rdb, queueName, numJobs)
	// dispatchJob1(rdb, queueName, 84055)
	//dispatchJob1(rdb, queueName, 84054)
}
