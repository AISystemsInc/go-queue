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
	// this is were you should run your processing tasks
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
		stop    = make(chan struct{}) // to command the program to stop creating jobs
		stopped = make(chan struct{}) // to notify us that the result routing has finished processing jobs
	)

	// start a new queue specifying the number of threads to use
	queue.Start(4, jobs, results)

	// start a goroutine to add jobs continuously
	// you could also pass the jobs channel to other parts of your system
	go func(stop chan struct{}) {
		var i int

		for {
			jobs <- &MyJob{Id: i}
			i++

			// checks if the stop signal has been sent
			// otherwise waits for 250ms
			select {
			case <-stop:
				close(jobs)
				return
			case <-time.After(time.Millisecond * 250):
			}
		}
	}(stop)

	// watch for results
	go func() {
		for result := range results {
			log.Printf("done: %d", result.(*MyJobResult).Data)
		}

		stopped <- struct{}{}
		close(stopped)
	}()

	// keep program alive, but imagine this was a http.listen or a very long process
	log.Println("keeping alive...")
	<-time.After(time.Second * 10)

	// send stop signal to job producer
	log.Println("stopping job producer")
	stop <- struct{}{}
	close(stop)

	// wait for result routine to handle all remaining results
	log.Println("waiting for remaining jobs to process...")
	<-stopped

	log.Println("done!")
}
