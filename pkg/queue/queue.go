package queue

type queue struct {
	jobs    <-chan Job
	results chan<- Result
	closed  closed
	workers int
}

type Job interface {
	Run() Result
}

type Result interface {
	Err() error
}

type Jobs chan Job
type Results chan Result
type closed chan struct{}

// Start creates worker goroutines and starts waiting for jobs to come through.
// Close the jobs channel to stop the queue, all existing jobs will be processed.
// After all workers have finished, the results channel will be closed.
func Start(workers int, jobs <-chan Job, results chan<- Result) {
	q := queue{
		workers: workers,
		jobs:    jobs,
		results: results,
		closed:  make(closed),
	}

	for i := 0; i < workers; i++ {
		go worker(q.jobs, q.results, q.closed)
	}

	// Handles worker close events
	// When all workers are closed we finally close the results channel
	go func() {
		var numClosed = 0
	closedLoop:
		for {
			select {
			case _, openOrMore := <-q.closed:
				if !openOrMore {
					break closedLoop
				}

				numClosed++

				if numClosed == q.workers {
					close(q.closed)
				}
			}
		}

		// we know we aren't going to receive anymore results now
		close(q.results)
	}()
}

// worker Reads the jobs channel and executes jobs it receives.
// It will only run one job at a time.
// The number of workers is equal to the number of threads requested
// at initialization.
func worker(jobs <-chan Job, results chan<- Result, closed closed) {
jobLoop:
	for {
		select {
		case job, openOrMore := <-jobs:
			if !openOrMore {
				closed <- struct{}{}
				break jobLoop
			}

			results <- job.Run()
		}
	}
}
