package my_modules

import "sync"

type CustomEvent struct {
	event_id      string
	subscribers_q chan *func(data interface{}) //incase if OnEvent called parallelly
	subscribers   []*func(data interface{})
}

func NewCustomEvent(event_id string) CustomEvent {
	e := &CustomEvent{
		event_id:      event_id,
		subscribers_q: make(chan *func(data interface{}),10),
	}
	return *e
}

func (e *CustomEvent) OnEvent(callback *func(data interface{})) {
	e.subscribers_q <- callback
}

func (e *CustomEvent) Emit(data interface{}) {
	for len(e.subscribers_q) > 0 {
		e.subscribers = append(e.subscribers, <-e.subscribers_q)
	}
	var wg sync.WaitGroup
	for _, subscriber := range e.subscribers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			(*subscriber)(data)
		}()
	}
	wg.Wait()
}

func (e *CustomEvent) Dispose() {
	close(e.subscribers_q)
	*e = CustomEvent{}
	e = nil
}
