/*
 * Example to demonstrate using channels and implement a reader and 
 * writer using channels.
 */
package main

import (
    "fmt"
    "context"
    "time"
    "log"
    "strings"
    "io"
)

var schedule_time int = 5
const MAX_LEN = 1500
const QUEUE_SIZE = 100

// Context created using WithCancel, WithDeadline, WithTimeout, or WithValue.
// When a Context is canceled, all Contexts derived from it are also canceled.
func Writer(ctx context.Context, out chan<- interface{}) error {

    var i int = 0
    for {
        content := fmt.Sprintf("hello world %d", i)
        buf := make([]byte, MAX_LEN)
        r := strings.NewReader(content)
        n, err := r.Read(buf)
        if err == io.EOF || err != nil {
           return err
        }

        v := buf[:n]

        select {
            case <-ctx.Done():
                fmt.Println("--cancelled")
                return ctx.Err()
            case out <- v:  // comes out of select loop when write to the channel is successfull
        }

        i++
    }
}

func ScheduleAlarm(alarm_time int, callback func() ()) (alarm_chan chan string) {
    alarm_chan = make(chan string)
    go func() {
        // Setting alarm.
        time.AfterFunc(time.Duration(alarm_time) * time.Second, func() {
            callback()
            alarm_chan <- "finished recording"
            close(alarm_chan)
        })
    }()
    return
}

func main() {
    ctx := context.Background()
    ctx, cancel := context.WithCancel(ctx)
    err_chan := make(chan error, 1)

    _ = ScheduleAlarm(schedule_time, func() {
        fmt.Println("alarm received")
        cancel()
    })

    msg_channel := make(chan interface{}, QUEUE_SIZE)

    // writer loop
    go Writer(ctx, msg_channel)

    // reader loop
    go func() {
        for {
            v := <-msg_channel
            switch v.(type) {
            case string:
               fmt.Println("got string value: %s", v)
            case []byte:
               fmt.Printf("got byte value: %q\n", v)
            default:
               fmt.Printf("got unknown value: %v", v)
            }
        }
    } ()

    select {
    case <-ctx.Done():
        fmt.Println("cancelled")
    case err := <-err_chan:
        if err != nil {
           fmt.Println("err-ed out")
           log.Fatal(err)
        }
    }
}

