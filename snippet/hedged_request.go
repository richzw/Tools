package snippet

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

var inflight int64

const BackupLimit = 10000

// call represents an active RPC.
type call struct {
	Name  string
	Reply interface{} // The reply from the function (*struct).
	Error error       // After completion, the error status.
	Done  chan *call  // Strobes when call is complete.
}

func (call *call) done() {
	select {
	case call.Done <- call:
	default:
		fmt.Printf("rpc: discarding Call reply due to insufficient Done chan capacity")
	}
}

func BackupRequest(backupTimeout time.Duration, fn func() (interface{}, error)) (interface{}, error) {
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	callCh := make(chan *call, 2)
	call1 := &call{Done: callCh, Name: "first"}
	call2 := &call{Done: callCh, Name: "second"}

	go func(c *call) {
		defer helpers.PanicRecover()
		c.Reply, c.Error = fn()
		c.done()
	}(call1)

	t := time.NewTimer(backupTimeout)
	select {
	case <-ctx.Done(): // cancel by context
		return nil, ctx.Err()
	case c := <-callCh:
		t.Stop()
		return c.Reply, c.Error
	case <-t.C:
		go func(c *call) {
			defer helpers.PanicRecover()
			defer atomic.AddInt64(&inflight, -1)
			if atomic.AddInt64(&inflight, 1) > BackupLimit {
				return
			}

			c.Reply, c.Error = fn()
			c.done()
		}(call2)
	}

	select {
	case <-ctx.Done(): // cancel by context
		return nil, ctx.Err()
	case c := <-callCh:
		return c.Reply, c.Error
	}
}
