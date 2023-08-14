# Jobber
Jobber is a lightweight Go library for running and managing background jobs. It provides an easy way to run, monitor, and control background tasks in a non-blocking fashion.

## Features
- Non-blocking job execution.
- Job state management using atomic operations.
- Error notifications with optional notification timeout.
- Flexible job implementation with custom timers and error handling.

## Installation
Install Jobber using go get:
```bash
go get -u github.com/go-pnp/jobber
```


## Usage
1) Define a job by implementing the Job interface:
```go
type MyJob struct{}

func (j *MyJob) Handle(ctx context.Context) error {
    // Your job logic here
    return nil
}

func (j *MyJob) Timer() *time.Timer {
    return time.NewTimer(5 * time.Second)
}

func (j *MyJob) ResetTimer(timer *time.Timer) {
    timer.Reset(5 * time.Second)
}
```
2) Initialize the Runner:
```go
jobInstance := &MyJob{}
runner := jobber.NewRunner(jobInstance)
```
3) Start and monitor the Runner:
```go
ctx := context.Background()
go func(){
    if err := runner.Start(ctx); err != nil {
        log.Fatal(err)
    }
}

// Monitor job errors (optional)
go func() {
    for err := range runner.Errors() {
        log.Println("Job error:", err)
    }
}()

// Stop the runner when done
if err := runner.Close(); err != nil {
    log.Fatal(err)
}
```


## Predefined Jobs
Jobber comes with a few predefined jobs that can be used out of the box:
- **IntervalJob**: Runs a job at a fixed interval.
- **InfinityJob**: Runs a job with no interval.
- **CronJob***: Runs a job at cron schedule

### IntervalJob
```go 
// Create a new IntervalJob that runs every 5 seconds
job := jobber.NewIntervalJob(
	true, // First iteration starts immediately
	5 * time.Second, // Execute every 5 seconds
	func(ctx context.Context) error {
		// Your job logic here
        return nil
    },
)
```

### CronJob
```go 
// Create a new IntervalJob that runs every 5 seconds
job, err := jobber.NewCronJob(
	true, // First iteration starts immediately
	"*/5 * * * *", // Execute every 5 minutes
	func(ctx context.Context) error {
		// Your job logic here
        return nil
    },
)
```

### InfinityJob
```go
// Create a new InfinityJob that runs forever
job := jobber.InfinityJob(func(ctx context.Context) error {
        // Your job logic here
        return nil
})
```
