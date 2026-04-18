package main

import (
    "time"
    "fmt"
)

func main() {
    fmt.Println("time:", time.Now())
    loc, err := time.LoadLocation("Asia/Shanghai")
    if err != nil {
        panic(err)
    }
    
    time.Local = loc

    fmt.Println("time:", time.Now())
}
