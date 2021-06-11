go-queue
========

An in-process job queue that utilizes all given threads continuously.


### Why is this better than sync.WaitGroup?

WaitGroup is a useful tool, but I've seen many cases where they are used inefficiently.

Lets propose a scenario:

I have `n` jobs which I have allocated `t` cores (and it should use no more!).

if i use a `sync.WaitGroup` i might write something like this:

```go
package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	var (
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
}

```

So what's the problem with that? It'll do that work but not use more than `4` threads.

Right! However, it will also underutilize the full capacity.

This is because when we call `wg.Wait()` we are waiting for **ALL** pending jobs to complete.
Essentially we are waiting for the slowest job to finish, of the `4`. So we could 
have up to `3` threads doing nothing at all!

This queue solves the problem by allowing each thread to take new work directly 
when it is available. Crucially never more than the number of threads allocated.

See this [example](examples/waitgroup_example/main.go) for a full code comparison.

Installation
------------

`go get github.com/AISystemsInc/go-queue`


Example usage
-------------

Checkout the [examples](examples) directory for usage modes.
