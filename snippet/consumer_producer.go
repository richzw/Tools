package snippet

import (
	"fmt"
	"sync"
	"time"
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

//https://mp.weixin.qq.com/s/hC_Hl4CQi725cirFwhxGfg

type SubWorkerNew struct {
	Id      int
	JobChan chan string
}

type W2New struct {
	SubWorkers []SubWorkerNew
	MaxNum     int
	ChPool     chan chan string
	QuitChan   chan struct{}
	Wg         *sync.WaitGroup
}

func NewW2(maxNum int) *W2New {
	chPool := make(chan chan string, maxNum)
	subWorkers := make([]SubWorkerNew, maxNum)
	for i := 0; i < maxNum; i++ {
		subWorkers[i] = SubWorkerNew{Id: i, JobChan: make(chan string)}
		chPool <- subWorkers[i].JobChan
	}
	wg := new(sync.WaitGroup)
	wg.Add(maxNum)

	return &W2New{
		MaxNum:     maxNum,
		SubWorkers: subWorkers,
		ChPool:     chPool,
		QuitChan:   make(chan struct{}),
		Wg:         wg,
	}
}

func (w *W2New) StartPool() {
	for i := 0; i < w.MaxNum; i++ {
		go func(wg *sync.WaitGroup, subWorker *SubWorkerNew) {
			defer wg.Done()
			for {
				select {
				case job := <-subWorker.JobChan:
					fmt.Printf("SubWorker %d processing job %s\n", subWorker.Id, job)
					time.Sleep(time.Second)
				case <-w.QuitChan:
					return
				}
			}
		}(w.Wg, &w.SubWorkers[i])
	}
}

func (w *W2New) Stop() {
	close(w.QuitChan)
	for i := 0; i < w.MaxNum; i++ {
		close(w.SubWorkers[i].JobChan)
	}
	w.Wg.Wait()
}

func (w *W2New) Dispatch(job string) {
	select {
	case jobChan := <-w.ChPool:
		jobChan <- job
	default:
		fmt.Println("All workers busy")
	}
}
func (w *W2New) AddWorker() {
	newWorker := SubWorkerNew{Id: w.MaxNum, JobChan: make(chan string)}
	w.SubWorkers = append(w.SubWorkers, newWorker)
	w.ChPool <- newWorker.JobChan
	w.MaxNum++
	w.Wg.Add(1)

	go func(subWorker *SubWorkerNew) {
		defer w.Wg.Done()

		for {
			select {
			case job := <-subWorker.JobChan:
				fmt.Printf("SubWorker %d processing job %s\n", subWorker.Id, job)
				time.Sleep(time.Second)
			case <-w.QuitChan:
				return
			}
		}
	}(&newWorker)
}

func (w *W2New) RemoveWorker() {
	if w.MaxNum > 1 {
		worker := w.SubWorkers[w.MaxNum-1]
		close(worker.JobChan)
		w.MaxNum--
		w.SubWorkers = w.SubWorkers[:w.MaxNum]
	}
}
