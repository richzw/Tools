package snippet

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type Broker struct {
	stopCh    chan struct{}
	publishCh chan interface{}
	subCh     chan chan interface{}
	unsubCh   chan chan interface{}
}

func NewBroker() *Broker {
	return &Broker{
		stopCh:    make(chan struct{}),
		publishCh: make(chan interface{}, 1),
		subCh:     make(chan chan interface{}, 1),
		unsubCh:   make(chan chan interface{}, 1),
	}
}

func (b *Broker) Start() {
	subs := map[chan interface{}]struct{}{}
	for {
		select {
		case <-b.stopCh:
			for msgCh := range subs {
				close(msgCh)
			}
			return
		case msgCh := <-b.subCh:
			subs[msgCh] = struct{}{}
		case msgCh := <-b.unsubCh:
			delete(subs, msgCh)
		case msg := <-b.publishCh:
			for msgCh := range subs {
				// msgCh is buffered, use non-blocking send to protect the broker:
				select {
				case msgCh <- msg:
				default:
				}
			}
		}
	}
}

func (b *Broker) Stop() {
	close(b.stopCh)
}

func (b *Broker) Subscribe() chan interface{} {
	msgCh := make(chan interface{}, 5)
	b.subCh <- msgCh
	return msgCh
}

func (b *Broker) Unsubscribe(msgCh chan interface{}) {
	b.unsubCh <- msgCh
}

func (b *Broker) Unsubscribe1(msgCh chan interface{}) {
	b.unsubCh <- msgCh
	close(msgCh)
}

func (b *Broker) Publish(msg interface{}) {
	b.publishCh <- msg
}

func main() {
	b := NewBroker()
	go b.Start()

	clientFunc := func(id int) {
		msgCh := b.Subscribe()
		for msg := range msgCh {
			fmt.Printf("Client %d got message: %v \n", id, msg)
		}
	}
	for i := 0; i < 3; i++ {
		go clientFunc(i)
	}

	go func() {
		for msgId := 0; ; msgId++ {
			b.Publish(fmt.Sprintf("msg# %d", msgId))
			time.Sleep(300 * time.Millisecond)
		}
	}()

	time.Sleep(time.Second)
}

//

type hub struct {
	sync.Mutex
	subs map[*subscriber]struct{}
}

func (h *hub) publish(ctx context.Context, msg *message) error {
	h.Lock()
	for s := range h.subs {
		s.publish(ctx, msg)
	}
	h.Unlock()

	return nil
}

func (h *hub) subscribe(ctx context.Context, s *subscriber) error {
	h.Lock()
	h.subs[s] = struct{}{}
	h.Unlock()

	go func() {
		select {
		case <-s.quit:
		case <-ctx.Done():
			h.Lock()
			delete(h.subs, s)
			h.Unlock()
		}
	}()

	go s.run(ctx)

	return nil
}

func (h *hub) unsubscribe(ctx context.Context, s *subscriber) error {
	h.Lock()
	delete(h.subs, s)
	h.Unlock()
	close(s.quit)
	return nil
}

func (h *hub) subscribers() int {
	h.Lock()
	c := len(h.subs)
	h.Unlock()
	return c
}

func newHub() *hub {
	return &hub{
		subs: map[*subscriber]struct{}{},
	}
}

type message struct {
	data []byte
}

type subscriber struct {
	sync.Mutex

	name    string
	handler chan *message
	quit    chan struct{}
}

func (s *subscriber) run(ctx context.Context) {
	for {
		select {
		case msg := <-s.handler:
			log.Println(s.name, string(msg.data))
		case <-s.quit:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (s *subscriber) publish(ctx context.Context, msg *message) {
	select {
	case <-ctx.Done():
		return
	case s.handler <- msg:
	default:
	}
}

func newSubscriber(name string) *subscriber {
	return &subscriber{
		name:    name,
		handler: make(chan *message, 100),
		quit:    make(chan struct{}),
	}
}
