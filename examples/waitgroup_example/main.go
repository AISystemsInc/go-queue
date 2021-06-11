package main

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/AISystemsInc/go-queue/pkg/queue"
)

// this example displays the issues with using sync.WaitGroup
// run this code and make note of the completion time of each.
// the queue will complete much faster because it is able to
// full utilize each thread.

func main() {
	waitGroup()
	queued()
}

func waitGroup() {
	var (
		start = time.Now()
		wg    sync.WaitGroup
		n     = 20
		t     = 4
	)

	for i := 1; i <= n; i++ {
		wg.Add(1)

		// start a job in a new routine.
		go func(idx int) {
			// simulate unpredictable job completion
			<-time.After(time.Second * time.Duration(1+rand.Intn(4)))
			log.Printf("done: %d", idx)
			wg.Done()
		}(i)

		// wait for last 4 to complete
		if i > 0 && i%t == 0 {
			log.Printf("syncing @%d", i)
			wg.Wait()
		}
	}

	log.Println("syncing last")
	wg.Wait()

	log.Printf("waitgroup took: %s", time.Since(start))
}

func queued() {
	var (
		start   = time.Now()
		jobs    = make(queue.Jobs)
		results = make(queue.Results)
		n       = 20
		t       = 4
	)

	queue.Start(t, jobs, results)

	go func() {
		for i := 1; i <= n; i++ {
			jobs <- &MyJob{
				Id: i,
			}
		}

		close(jobs)
	}()

	for result := range results {
		log.Printf("done: %d", result.(*MyJobResult).Data)
	}

	log.Printf("queue took: %s", time.Since(start))
}

//==============================================================================
// Define a job
//==============================================================================

type MyJob struct {
	Id int
}

func (m *MyJob) Run() queue.Result {
	// this is were you should run your processing tasks
	<-time.After(time.Second * time.Duration(rand.Intn(5)))

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
