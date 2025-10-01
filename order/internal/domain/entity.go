package domain

type Entity struct {
	events []Event
}

func (e *Entity) PullEvents() []Event {
	events := e.events
	e.events = []Event{}
	return events
}
