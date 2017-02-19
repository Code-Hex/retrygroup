package retrygroup

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"golang.org/x/net/context"
)

// A Group is a collection of goroutines working on subtasks that are part of
// the same overall task.
type Group struct {
	cancel func()

	w        io.Writer
	interval time.Duration
	wg       sync.WaitGroup
}

// WithContext returns a new Group and an associated Context derived from ctx.
func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel, w: os.Stderr}, ctx
}

func (g *Group) SetWriter(out io.Writer) {
	g.w = out
}

// Wait blocks until all function calls from the Go method have returned.
func (g *Group) Wait() {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
}

func (g *Group) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		for {
			if err := f(); err != nil {
				fmt.Fprint(g.w, err.Error())
			}
		}
	}()
}

func (g *Group) GoTimes(cnt int, f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		for i := 0; i < cnt; i++ {
			if err := f(); err != nil {
				fmt.Fprint(g.w, err.Error())
			}
		}
	}()
}
