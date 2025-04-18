package painter

import (
	"image"
	"sync"
	"time"

	"golang.org/x/exp/shiny/screen"
)

type Receiver interface {
	Update(t screen.Texture)
}

type Loop struct {
	Receiver Receiver

	next screen.Texture
	prev screen.Texture

	mq messageQueue

	State *State

	stop chan struct{}
	Done chan struct{}
	stopReq bool
	mu sync.Mutex
	loopRunning bool
}

var size = image.Pt(800, 800)

func (l *Loop) Start(s screen.Screen) {
	l.mu.Lock()
	if l.loopRunning {
		l.mu.Unlock()
		return
	}
	l.loopRunning = true
	l.mu.Unlock()

	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)
	l.mq = messageQueue{queue: make(chan Operation, 1000)}
	l.stop = make(chan struct{})
	l.Done = make(chan struct{})
	l.State = DefaultState()

	go l.run()
}

func (l *Loop) Post(op Operation) {
	l.mq.push(op)
}

func (l *Loop) StopAndWait() {
	l.mu.Lock()
	if !l.loopRunning {
		l.mu.Unlock()
		if l.Done != nil {
			<-l.Done
		}
		return
	}
	l.stopReq = true
	l.mu.Unlock()
	<-l.Done
}

func (l *Loop) run() {
	defer func() {
		l.mu.Lock()
		l.loopRunning = false
		l.mu.Unlock()
		close(l.Done)
		if l.next != nil {
			l.next.Release()
		}
		if l.prev != nil {
			l.prev.Release()
		}
	}()

	for {
		l.mu.Lock()
		if l.stopReq && l.mq.empty() {
			l.mu.Unlock()
			return
		}
		l.mu.Unlock()

		select {
		case op := <-l.mq.queue:
			needsUpdate := op.Do(l.next, l.State)

			if needsUpdate {
				DrawStateOp{}.Do(l.next, l.State)
				l.Receiver.Update(l.next)
				l.next, l.prev = l.prev, l.next
			}

		case <-time.After(time.Millisecond * 100):
		}
	}
}

type messageQueue struct {
	queue chan Operation
}

func (mq *messageQueue) push(op Operation) {
	mq.queue <- op
}

func (mq *messageQueue) pull() Operation {
	select {
	case op := <-mq.queue:
		return op
	default:
		return nil
	}
}

func (mq *messageQueue) empty() bool {
	return len(mq.queue) == 0
}
