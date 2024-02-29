package snippet

import (
	"fmt"
	"sync"
)

type Promise struct {
	wg  sync.WaitGroup
	res string
	err error
}

func NewPromise(f func() (string, error)) *Promise {
	p := &Promise{}
	p.wg.Add(1)
	go func() {
		p.res, p.err = f()
		p.wg.Done()
	}()
	return p
}

func (p *Promise) Then(r func(string), e func(error)) {
	go func() {
		p.wg.Wait()
		if p.err != nil {
			e(p.err)
			return
		}
		r(p.res)
	}()
}

// With channel
// http://www.home.hs-karlsruhe.de/~suma0002/publications/events-to-futures.pdf
type Comp struct {
	value interface{}
	ok    bool
}

type Future chan Comp

func future(f func() (interface{}, bool)) Future {
	future := make(chan Comp)

	go func() {
		v, o := f()
		c := Comp{v, o}
		for {
			future <- c
		}
	}()

	return future
}

type Promise struct {
	lock chan int
	ft   Future
	full bool
}

func promise() Promise {
	return Promise{make(chan int, 1), make(chan Comp), false}
}

func (pr Promise) future() Future {
	return pr.ft
}

//

// Promise structure
type PromiseV1 struct {
	result chan interface{}
	err    chan error
}

// NewPromise initializes and returns a new Promise.
func NewPromiseV1(f func() (interface{}, error)) *PromiseV1 {
	p := &PromiseV1{
		result: make(chan interface{}, 1),
		err:    make(chan error, 1),
	}

	go func() {
		res, err := f()
		if err != nil {
			p.err <- err
			return
		}
		p.result <- res
	}()

	return p
}

// Then method for handling success
func (p *PromiseV1) Then(successFunc func(interface{}) interface{}) *PromiseV1 {
	newPromise := &PromiseV1{
		result: make(chan interface{}, 1),
		err:    make(chan error, 1),
	}

	go func() {
		select {
		case res := <-p.result:
			newResult := successFunc(res)
			newPromise.result <- newResult
		case err := <-p.err:
			newPromise.err <- err
		}
	}()

	return newPromise
}

// Catch method for handling errors
func (p *PromiseV1) Catch(failFunc func(error)) {
	go func() {
		for err := range p.err {
			failFunc(err)
		}
	}()
}

// Example usage
func main() {
	p := NewPromiseV1(func() (interface{}, error) {
		// Simulate some work
		// Return either a result or an error
		return "Success result", nil
	})

	p.Then(func(data interface{}) interface{} {
		fmt.Printf("First Then: %v\n", data)
		return "Result from first Then"
	}).Then(func(data interface{}) interface{} {
		fmt.Printf("Second Then: %v\n", data)
		return nil
	}).Catch(func(err error) {
		fmt.Printf("Caught error: %v\n", err)
	})
}

func Clone[S ~[]E, E any](s S) S {
	return append(s[:0:0], s...)
}
