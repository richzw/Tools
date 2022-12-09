package snippet

import (
	"errors"
	"log"
	"sync"
)

type Message interface {
}

type Stage interface {
	Process(stage Message) ([]Message, error)
}

type PipelineOpts struct {
	Concurrency int
}
type Pipeline interface {
	AddPipe(pipe Stage, opt *PipelineOpts)
	Start() error
	Stop() error
	Input() chan<- Message
	Output() <-chan Message
}

var ErrConcurrentPipelineEmpty = errors.New("concurrent pipeline empty")

type ConcurrentPipeline struct {
	workerGroups []StageWorker
}

func (c *ConcurrentPipeline) AddPipe(pipe Stage, opt *PipelineOpts) {
	if opt == nil {
		opt = &PipelineOpts{Concurrency: 1}
	}

	var input = make(chan Message, 10)
	var output = make(chan Message, 10)

	for _, i := range c.workerGroups {
		input = i.Output()
	}

	worker := NewWorkerGroup(opt.Concurrency, pipe, input, output)
	c.workerGroups = append(c.workerGroups, worker)
}

func (c *ConcurrentPipeline) Output() <-chan Message {
	sz := len(c.workerGroups)
	return c.workerGroups[sz-1].Output()
}

func (c *ConcurrentPipeline) Input() chan<- Message {
	return c.workerGroups[0].Input()
}

func (c *ConcurrentPipeline) Start() error {
	if len(c.workerGroups) == 0 {
		return ErrConcurrentPipelineEmpty
	}

	for i := 0; i < len(c.workerGroups); i++ {
		g := c.workerGroups[i]
		g.Start()
	}

	return nil
}

func (c *ConcurrentPipeline) Stop() error {

	for _, i := range c.workerGroups {
		close(i.Input())
		i.WaitStop()
	}

	sz := len(c.workerGroups)
	close(c.workerGroups[sz-1].Output())
	return nil
}

func NewConcurrentPipeline() Pipeline {
	return &ConcurrentPipeline{}
}

type StageWorker struct {
	wg          *sync.WaitGroup
	input       chan Message
	output      chan Message
	concurrency int
	pipe        Stage
}

func (w *StageWorker) Start() error {
	for i := 0; i < w.concurrency; i++ {
		w.wg.Add(1)

		go func() {
			defer w.wg.Done()
			for i := range w.Input() {
				result, err := w.pipe.Process(i)
				if err != nil {
					log.Println(err)
					continue
				}
				for _, r := range result {
					w.Output() <- r
				}
			}
		}()
	}

	return nil

}

func (w *StageWorker) WaitStop() error {
	w.wg.Wait()
	return nil
}

func (w *StageWorker) Input() chan Message {
	return w.input
}

func (w *StageWorker) Output() chan Message {
	return w.output
}

func NewWorkerGroup(concurrency int, pipe Stage, input chan Message, output chan Message) StageWorker {

	return StageWorker{
		wg:          &sync.WaitGroup{},
		input:       input,
		output:      output,
		concurrency: concurrency,
		pipe:        pipe,
	}

}
