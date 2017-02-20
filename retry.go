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

	w   io.Writer
	wg  sync.WaitGroup
	ctx context.Context

	canbackoff bool
}

// WithContext returns a new Group and an associated Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first.
func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{
		cancel: cancel,
		w:      os.Stderr,
		ctx:    ctx,
	}, ctx
}

// EnableBackoff will enable exponential backoff.
func (g *Group) EnableBackoff() {
	g.canbackoff = true
}

// SetWriter can set the io.Writer when there is an error in the RetryGo.
func (g *Group) SetWriter(out io.Writer) {
	g.w = out
}

// Wait blocks until all function calls from the RetryGo method have returned.
func (g *Group) Wait() {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
}

// RetryGo calls the given function in a new goroutine.
//
// If the value of the argument representing the number of times is
// set to `cnt <= 0`, retry of that routine will be done forever.
// Other than that, it will execute it for that number of times.
func (g *Group) RetryGo(cnt int, f func(int) error) {
	g.wg.Add(1)

	times := uint(cnt)

	go func() {
		defer g.wg.Done()
		for i := uint(0); i < times || cnt <= 0; i++ {
			select {
			case <-g.ctx.Done():
				return
			default:
				if err := f(int(i + 1)); err != nil {
					fmt.Fprintln(g.w, err.Error())
				} else {
					break
				}
			}
			// Exponential backoff
			if g.canbackoff {
				time.Sleep(time.Second * time.Duration(1<<i))
			}
		}
	}()
}
