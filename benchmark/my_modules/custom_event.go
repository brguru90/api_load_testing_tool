package my_modules

type CustomEvent struct {
	event_id      string
	subscribers_q chan *func(data interface{})
	subscribers   []*func(data interface{})
}

func NewCustomEvent(event_id string) CustomEvent {
	e := &CustomEvent{
		event_id:      event_id,
		subscribers_q: make(chan *func(data interface{}), 1000),
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
	for _, subscriber := range e.subscribers {
		(*subscriber)(data)
	}
}

func (e *CustomEvent) Dispose() {
	close(e.subscribers_q)
	*e = CustomEvent{}
	e = nil
}
