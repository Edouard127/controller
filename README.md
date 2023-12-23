# Tor Controller

Simple Tor controller for Golang 1.21+.

## Usage

```go get github.com/Edouard127/controller```

```go
package main

import (
    "fmt"
    "github.com/Edouard127/controller"
)

func main() {
    c, err := controller.NewController("127.0.0.1:9051")
    if err != nil {
        panic(err)
    }

    err = c.Authenticate("") // Or c.Authenticate("password") if you have one
    if err != nil {
        panic(err)
    }

    fmt.Println(c.GetInfo("version"))
}
```