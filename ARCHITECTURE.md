Architecture
============

`queue` is a simple job processor that runs entirely in-memory and, as such, has no persistent state (db).

`queue` is different from the standard `sync.WaitGroup`. It can keep the workers saturated with jobs while there are >= jobs to workers.
Using the alternative `sync.WaitGroup` loop there would be a decline in thread usage as the current iteration's jobs decreased. Leaving potential processing power underutilised.


Jobs are executed by a predefined number of workers, each of which can only process one job at a time.

As a worker is finished with a job, it immediately takes a new one from the jobs channel. Results are pushed into the results channel.

If there is a block on the receiving end of the results channel then the queue will pause, although active workers will finish processing their current job.

This process continues until the jobs channel is closed.

### Diagram of the architecture:

```
                        ╔═══════════════════════════════════════════════════════╗
 ┌──────────────────┐   ║Queue (N Workers)       ╔════════════════════════════╗ ║
 │       Jobs       │──▶║                        ║Worker                      ║ ║
 └──────────────────┘   ║┌──────────────────┐    ║                            ║ ║
 ┌──────────────────┐   ║│ Create N workers │─┐  ║┌──────────────────────────┐║ ║
 │     Results      │──▶║└──────────────────┘ │  ║│          <-Jobs          │║ ║
 └──────────────────┘   ║┌──────────────────┐ │  ║└──────────────────────────┘║ ║
                        ║│     <-Closed     │ │  ║┌──────────────────────────┐║ ║
  ┌──────────────────┐  ║└──────────────────┘ ├─▶║│  Closed<- (Jobs closed)  │║ ║
  │Each worker can   │▒ ║┌──────────────────┐ │  ║└──────────────────────────┘║ ║
  │only accept one   │▒ ║│ if N workers are │ │  ║┌──────────────────────────┐║ ║
  │job at a time     │▒ ║│      closed      │ │  ║│  Results<- (Job.Start())   │║ ║
  └──────────────────┘▒ ║│  Close(Closed &  │ │  ║└──────────────────────────┘║ ║
   ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒ ║│     Results)     │ │  ║                            ║ ║
                        ║└──────────────────┘ │  ╚════════════════════════════╝ ║
  ┌──────────────────┐  ║                     │  ╔════════════════════════════╗ ║
  │Close Jobs to     │▒ ║                     │  ║Worker                      ║ ║
  │safely shutdown   │▒ ║                     │  ║                            ║ ║
  │queue             │▒ ║                     │  ║┌──────────────────────────┐║ ║
  └──────────────────┘▒ ║                     │  ║│          <-Jobs          │║ ║
   ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒ ║                     │  ║└──────────────────────────┘║ ║
                        ║                     │  ║┌──────────────────────────┐║ ║
  ┌──────────────────┐  ║                     └─▶║│  Closed<- (Jobs closed)  │║ ║
  │Keep reading on   │▒ ║                        ║└──────────────────────────┘║ ║
  │Results after     │▒ ║                        ║┌──────────────────────────┐║ ║
  │closing Jobs to   │▒ ║                        ║│  Results<- (Job.Start())   │║ ║
  │read all results  │▒ ║                        ║└──────────────────────────┘║ ║
  └──────────────────┘▒ ║                        ║                            ║ ║
   ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒ ║                        ╚════════════════════════════╝ ║
                        ╚═══════════════════════════════════════════════════════╝
```
