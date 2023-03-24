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

// without waitgroup

const producerCount int = 3
const consumerCount int = 3

var workers []*producers

type producers struct {
	myQ  chan string
	quit chan bool
	id   int
}

func execute(jobQ chan<- string, workerPool chan *producers, allDone chan<- bool) {
	for _, j := range messages {
		jobQ <- j
	}
	close(jobQ)
	for _, w := range workers {
		w.quit <- true
	}
	close(workerPool)
	allDone <- true
}

func produce(jobQ <-chan string, p *producers, workerPool chan *producers) {
	for {
		select {
		case msg := <-jobQ:
			{
				workerPool <- p
				if len(msg) > 0 {
					fmt.Printf("Job \"%v\" produced by worker %v\n", msg, p.id)
				}
				p.myQ <- msg
			}
		case <-p.quit:
			return
		}
	}
}

func consume(cIdx int, workerPool <-chan *producers) {
	for {
		worker := <-workerPool
		if msg, ok := <-worker.myQ; ok {
			if len(msg) > 0 {
				fmt.Printf("Message \"%v\" is consumed by consumer %v from worker %v\n", msg, cIdx, worker.id)
			}
		}
	}
}

func test() {
	jobQ := make(chan string)
	allDone := make(chan bool)
	workerPool := make(chan *producers)

	for i := 0; i < producerCount; i++ {
		workers = append(workers, &producers{
			myQ:  make(chan string),
			quit: make(chan bool),
			id:   i,
		})
		go produce(jobQ, workers[i], workerPool)
	}

	go execute(jobQ, workerPool, allDone)

	for i := 0; i < consumerCount; i++ {
		go consume(i, workerPool)
	}
	<-allDone
}
