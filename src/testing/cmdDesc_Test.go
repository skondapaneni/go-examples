package main

import (
   "cli"
   "os"
    "log"
    "bufio"
    "fmt"
)


func main() {
    cmd_desc := cli.NewCmd()
    cmd_desc.Init()

    file, err := os.Open("file.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    // Default scanner is bufio.ScanLines. Lets use ScanWords.
    // Could also use a custom function of SplitFunc type
    // scanner.Split(bufio.ScanWords)

    // Scan for next token. 
    success := true
    for success {
        success = scanner.Scan() 
        if (success) {
            fmt.Println("line = " + scanner.Text())  // token in unicode-char
            fmt.Printf("first_list = %+v\n", cmd_desc.Complete(scanner.Text()) )
        }
    }

    if success == false {
        // False on error or EOF. Check error
        err = scanner.Err()
        if err == nil {
            log.Println("Scan completed and reached EOF")
        } else {
            log.Fatal(err)
        }
    }
}
