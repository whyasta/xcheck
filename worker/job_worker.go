package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// JobCounter holds the count of jobs processed and provides synchronization across workers
type JobCounter struct {
	mu    sync.Mutex
	count int
}

// Worker listens to the Redis queue and processes jobs
func worker(id int, rdb *redis.Client, queueName string, wg *sync.WaitGroup, counter *JobCounter) {
	defer wg.Done()

	for {
		// Fetch job from the Redis queue (blocking pop)
		job, err := rdb.BLPop(ctx, 0*time.Second, queueName).Result()
		if err != nil {
			log.Printf("Worker %d: Error retrieving job: %v\n", id, err)
			continue
		}

		// job[1] is the actual job data because BLPop returns a slice with the key and value
		fmt.Printf("Worker %d started job: %s\n", id, job[1])

		// Simulate processing time (adjust this based on actual job processing)
		processJob(id, job[1])

		fmt.Printf("Worker %d completed job: %s\n", id, job[1])

		// Increment the job counter and check if we reached a threshold of 100
		counter.incrementAndNotify()
	}
}

// processJob simulates processing a job, replace this with actual job logic
func processJob(id int, jobData string) {
	//time.Sleep(1 * time.Second) // Simulate work with sleep
	start := time.Now()

	// Simulate an HTTP GET request
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		fmt.Printf("Worker %d job %s: Error making request: %v\n", id, jobData, err)
		return
	}
	defer resp.Body.Close()

	// body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Worker %d job %s: Request completed in %v\n", id, jobData, time.Since(start))
	// fmt.Printf("Worker %d job %s: Response: %s\n", id, jobData, body[:100]) // Print first 100 bytes of the response
}

// Increment job counter and notify every 100 jobs
func (jc *JobCounter) incrementAndNotify() {
	jc.mu.Lock()
	defer jc.mu.Unlock()

	jc.count++
	if jc.count%5 == 0 {
		// Notification can be logging, or you could trigger an external service here (email, webhook, etc.)
		fmt.Printf("Notification: %d jobs have been processed.\n", jc.count)
	}
}

func main() {
	// Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
		DB:   0,                // Default DB
	})

	// Redis queue name
	queueName := "jobQueue"
	numWorkers := 50

	// Job counter to track total processed jobs across workers
	counter := &JobCounter{}

	// WaitGroup to ensure main function waits for workers
	var wg sync.WaitGroup

	// Start workers to listen to the Redis queue
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, rdb, queueName, &wg, counter)
	}

	fmt.Println("Start")

	// Wait for all workers (they run indefinitely in this case)
	wg.Wait()
}
