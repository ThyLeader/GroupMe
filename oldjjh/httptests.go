package main

import (
    "net/http"
    "fmt"
    "io/ioutil"
    "net/url"
    "time"
)

func main() {
    for {
        resp, err := http.PostForm("https://api.groupme.com/v3/bots/post",
            url.Values{"text" : {"DONT MESS WITH ME"}, "bot_id" : {"a9366914d71a5fb96d0408dd42"}})

        if err != nil {
            fmt.Println("error happened while getting the response", err)
            return
        }

        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)

        if err != nil {
            fmt.Println("error happened while reading the body", err)
            return
        }

        fmt.Println(string(body[:]))
        time.Sleep(100 * time.Millisecond)
    }
}