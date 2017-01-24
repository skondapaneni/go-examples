package libinfra

import (
)

/*
 * A worker pool is a re-usable library for creating  a pool
 * of threads to work on a job in jobQueue.
 *
 * All the caller does is :
 * 1. create a job queue channel
 *    var JobQueue chan libinfra.Job = make(chan libinfra.Job)
 *
 * 2. Create a AppJob of type Job which adheres to the Job interface.
 *     See CmdJob example below
 *
 * 3. Create a jobHandler which will run the job
 *    func jobHandler(job AppJob) {
 *    }
 *
 * 4. Create a dispatcher
 *    dispatcher := NewDispatcher(MAX_WORKERS, jobQueue, jobHandler)
 *
 * 5. And run the dispacther as given below.
 *    dispatcher.Run()
 *
 * 6. At runtime when work needs to be dipatched ..
 *        create and initialize a new appJob
 *      work := AppJob{app job data ..}
 *
 * 7.  - push the work onto the queue
 *    dispatcher.DispatchJobToWorker(work)
 *
 * // My workerCallback which will do the task
 * func workerCallback(job Job) {
 *      cmds := job.GetJobData()
 *      switch v := cmds.(type) {
 *        case []string:
 *        for _, cmd := range v {
 *            fmt.Printf("\nworker cmd " + cmd)
 *        }
 *      default:
 *        fmt.Println(cmds, "is of a type I don't know how to handle")
 *      }
 * }

 * func main() {
 *
 *    const MAX_WORKERS = 5
 *    var jobQueue chan libinfra.Job = make(chan libinfra.Job)
 *    var done = make(chan bool)
 *
 *    dispatcher := NewDispatcher(MAX_WORKERS, jobQueue, workerCallback)
 *    dispatcher.Run()
 *
 *    cmdJob1 := NewCmdJob()
 *    cmdJob1.AddJobData("one")
 *    cmdJob1.AddJobData("two")
 *
 *    cmdJob2 := NewCmdJob()
 *    cmdJob2.AddJobData("1")
 *    cmdJob2.AddJobData("2")
 *
 *    dispatcher.DispatchJobToWorker(cmdJob1)
 *    dispatcher.DispatchJobToWorker(cmdJob2)
 *
 *    // Main thread to keep spinning..
 *    for {
 *        select {
 *        case <-done:
 *            fmt.Printf("app terminated unexpectedly")
 *        default:
 *            break
 *        } // end of select
 *    }
 *
 *
 * }
 */

/**
 * Job defines an interface for all job types that can be dispatched to worker.
 */
type Job interface {
    SetJobTag(tag string)
    GetJobTag() string
    AddJobData(interface{}) error
    GetJobData() interface{} // job or task data
    IsJobEmpty() bool
}


/**
 * A handler method type for the job Handler which handles the job
 */
type JobHandler func(job Job)



/**
 * Worker represents the worker that executes the job
 */
type Worker struct {
    WorkerPool chan chan Job
    JobChannel chan Job /* A buffered channel that we can send work requests on. */
    quit       chan bool
    jobHandler JobHandler
}

func NewWorker(workerPool chan chan Job, jobHandler JobHandler) Worker {
    return Worker{
        WorkerPool: workerPool,
        JobChannel: make(chan Job),
        quit:       make(chan bool),
        jobHandler: jobHandler,
    }
}

/**
 * Start method starts the run loop for the worker, listening for a quit channel in
 * case we need to stop it
 */
func (w Worker) Start() {
    go func() {
        for {
            // register the current worker into the worker queue.
            w.WorkerPool <- w.JobChannel

            select {
            case job := <-w.JobChannel:
                w.jobHandler(job)

            case <-w.quit:
                // we have received a signal to stop
                return
            }
        }
    }()
}

/**
 *  Stop signals the worker to stop listening for work requests.
 */
func (w Worker) Stop() {
    go func() {
        w.quit <- true
    }()
}

/**
 * Dispatcher dispatches the jobs to available worker
 */
type Dispatcher struct {
    // A pool of workers channels that are registered with the dispatcher
    WorkerPool chan chan Job
    maxWorkers int
    handler    JobHandler
    jobQueue   chan Job
}

func NewDispatcher(maxWorkers int, jobQueue chan Job,
    handler JobHandler) *Dispatcher {
    pool := make(chan chan Job, maxWorkers)
    return &Dispatcher{WorkerPool: pool,
        maxWorkers: maxWorkers,
        handler:    handler,
        jobQueue:   jobQueue,
    }
}

/**
 * Initializes the Worker threads
 */
func (d *Dispatcher) Run() {
    // starting n number of workers
    for i := 0; i < d.maxWorkers; i++ {
        worker := NewWorker(d.WorkerPool, d.handler)
        worker.Start()
    }
    go d.dispatch()
}

/**
 * The dispatcher loop
 */
func (d *Dispatcher) dispatch() {
    for {
        select {
        case job := <-d.jobQueue:
            // a job request has been received
            go func(job Job) {
                // try to obtain a worker job channel that is available.
                // this will block until a worker is idle
                jobChannel := <-d.WorkerPool

                // dispatch the job to the worker job channel
                jobChannel <- job
            }(job)
        }
    }
}

func (d *Dispatcher) DispatchJobToWorker(job Job) {
    d.jobQueue <- job
}

