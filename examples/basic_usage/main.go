package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/AISystemsInc/go-queue/pkg/queue"
)

//==============================================================================
// Define a job
//==============================================================================

type MyJob struct {
	Id int
}

func (m *MyJob) Run() queue.Result {
	// this is were you should run your processing tasks.
	<-time.After(time.Second * time.Duration(1+rand.Intn(5)))

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
// Main program
//==============================================================================

func main() {
	// create our channels
	// depending on your requirements you can make these buffered
	var (
		jobs    = make(queue.Jobs)
		results = make(queue.Results)
	)

	// start a new queue specifying the number of threads to use
	queue.Start(4, jobs, results)

	// start a goroutine to add jobs
	go func() {
		for i := 0; i < 16; i++ {
			jobs <- &MyJob{Id: i}
		}

		close(jobs)
	}()

	// watch for results
	for result := range results {
		log.Printf("done: %d", result.(*MyJobResult).Data)
	}
}
