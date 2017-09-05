
package main

import (
   "fmt"
   "libinfra"
)

// Example CmdJob - A type of Job which implements the Job Interface
type CmdJob struct {
    cmds []string
    tag  string
}

func NewCmdJob(tag string) *CmdJob {
    return &CmdJob{
        cmds: make([]string, 0),
        tag: tag,
    }
}

func (self *CmdJob) SetJobTag(tag string) {
    self.tag = tag
}

func (self *CmdJob) GetJobTag() string {
    return self.tag
}

func (self *CmdJob) GetJobData() interface{} {
    return self.cmds
}

func (self *CmdJob) AddJobData(cmd interface{}) error {
    switch v := cmd.(type) {
    case string:
        self.cmds = append(self.cmds, v)

    default:
        fmt.Printf("%T is a type which is not supported", v)
    }
    return nil
}

func (self *CmdJob) IsJobEmpty() bool {
    return (len(self.cmds) == 0)
}

// My workerCallback which will do the task
func workerCallback(job libinfra.Job) {
    cmds := job.GetJobData()
    switch v := cmds.(type) {
    case []string:
        for _, cmd := range v {
            fmt.Printf("\nworker cmd %s\n",  cmd)
        }
    default:
        fmt.Println(cmds, "is of a type I don't know how to handle")
    }
}


func main() {

    const MAX_WORKERS = 5
    var jobQueue chan libinfra.Job = make(chan libinfra.Job)
    var done = make(chan bool)

    dispatcher := libinfra.NewDispatcher(MAX_WORKERS, jobQueue, workerCallback)
    dispatcher.Run()

    cmdJob1 := NewCmdJob("cmdJob1")
    cmdJob1.AddJobData("one")
    cmdJob1.AddJobData("two")

    cmdJob2 := NewCmdJob("cmdJob2")
    cmdJob2.AddJobData("1")
    cmdJob2.AddJobData("2")

    dispatcher.DispatchJobToWorker(cmdJob1)
    dispatcher.DispatchJobToWorker(cmdJob2)

    // Main thread to keep spinning..
    for {
        select {
        case <-done:
            fmt.Printf("app terminated unexpectedly")
        default:
            break
        } // end of select
    }

    fmt.Printf("app terminated unexpectedly 2")
}
