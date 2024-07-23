package processors

import (
	"fmt"

	"github.com/gocraft/work"
)

func TestJob(job *work.Job) error {
	fmt.Println("test job TestJob")
	// Extract arguments:
	addr := job.ArgString("email_address")
	if err := job.ArgError(); err != nil {
		return err
	}

	fmt.Println("Testing job", addr)
	return nil
}
