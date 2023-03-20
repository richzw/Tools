package snippet

import (
	"fmt"
	"sync"
)

// multi producer and multi consumer
func producer(ch chan<- any, src chan any, wg *sync.WaitGroup) {
	defer wg.Done()

	for v := range src {
		// produce job
		ch <- v
	}
}

func consumer(ch <-chan any, wg *sync.WaitGroup) {
	defer wg.Done()

	for v := range ch {
		// consumer job
		fmt.Println(v)
	}
}

func prodCons() {
	ch := make(chan any, 10)
	source := make(chan any, 5)

	wp := &sync.WaitGroup{}
	wc := &sync.WaitGroup{}

	for i := 0; i < 5; i++ {
		wp.Add(1)
		go producer(ch, source, wp)
	}

	for i := 0; i < 8; i++ {
		wc.Add(1)
		go consumer(ch, wc)
	}

	for i := 0; i < 10000; i++ {
		// simulate job insert
		source <- i
	}

	wp.Wait()
	close(ch)
	wc.Wait()
}
