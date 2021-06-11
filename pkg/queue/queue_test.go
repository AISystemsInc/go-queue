package queue_test

import (
	"testing"
	"time"

	"github.com/AISystemsInc/go-queue/pkg/queue"
)

//==============================================================================
// Define a job
//==============================================================================

type TimeoutJob struct {
	Id      int
	Timeout time.Duration
}

func (m *TimeoutJob) Run() queue.Result {
	<-time.After(m.Timeout)

	return &MyJobResult{
		err:  nil,
		Data: m.Id,
	}
}

//==============================================================================
// Define the job result type
//==============================================================================

type MyJobResult struct {
	err  error
	Data int
}

func (m *MyJobResult) Err() error {
	return m.err
}

//==============================================================================
// Tests
//==============================================================================

func TestQueue(t *testing.T) {
	t.Run("works concurrently", func(t *testing.T) {
		var (
			jobs    = make(queue.Jobs)
			results = make(queue.Results)
			numJobs = 20
			threads = 4
		)

		queue.Start(threads, jobs, results)

		go func() {
			for i := 0; i < numJobs; i++ {
				jobs <- &TimeoutJob{Timeout: time.Second}
			}
			close(jobs)
		}()

		var start = time.Now()

		for range results {
		}

		var timeTaken = time.Since(start)

		if timeTaken-(time.Millisecond*20) > time.Second*5 {
			t.Errorf("queue did not complete on time, took %s, expected: %s", timeTaken, time.Second*time.Duration(numJobs/threads))
		}
	})
}
