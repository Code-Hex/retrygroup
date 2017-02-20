retrygroup
----------
[![GoDoc](https://godoc.org/github.com/Code-Hex/retrygroup?status.svg)](https://godoc.org/github.com/Code-Hex/retrygroup)  
Package retrygroup provides synchronization, Context cancelation for groups of retry goroutines working on subtasks of a common task.

# Synopsis
```go
package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Code-Hex/retrygroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g, _ := retrygroup.WithContext(ctx)
	g.EnableBackoff()

	go func() {
		<-time.After(16 * time.Second)
		if cancel != nil {
			fmt.Println("Finish!!")
			cancel()
		}
	}()

	g.RetryGo(3, func(i int) error {
		fmt.Printf("Hello: %d\n", i)
		return errors.New("Try error")
	})

	g.RetryGo(-1, func(i int) error {
		fmt.Println("Never!!")
		return errors.New("Try never error")
	})

	g.Wait()
}
```
See [eg](https://github.com/Code-Hex/retrygroup/tree/master/eg)
