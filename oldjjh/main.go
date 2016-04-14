package main

import (
    "fmt"
    "time"
)

func main() {
    z := 0
    for {
        for i := "Connecting to proxy";; i += "." {
            go connectToProxy()
            fmt.Printf("\r%s", i)
            time.Sleep(500 * time.Millisecond)
            z++
            if z > 5 {
                fmt.Println("")
                break
            }
        }
        z = 0
        for i := "Loading website";; i += "." {
            go loadWebsite()
            fmt.Printf("\r%s", i)
            time.Sleep(500 * time.Millisecond)
            z++
            if z > 3 {
                fmt.Println("")
                break
            }
        }
        z = 0
        for i := "Sending vote";; i += "." {
            go sendVote()
            fmt.Printf("\r%s", i)
            time.Sleep(500 * time.Millisecond)
            z++
            if z > 2 {
                fmt.Println("")
                break
            }
        }
        fmt.Println("Vote sent. Resetting")
        go resetProxy()
        time.Sleep(500 * time.Millisecond)
        fmt.Println("Reset. Loading new proxy.")
        time.Sleep(3000 * time.Millisecond)
    }
    z = 0
}

func sendVote() {
    // beep boop send vote
}

func resetProxy() {
    //beep boop reset proxy from list
}

func connectToProxy() {
    //beep boop connect to proxy
}

func loadWebsite() {
    //beep boop load site
}
