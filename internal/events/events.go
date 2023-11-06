package events

type Fetcher interface {
	Fetch(int) ([]Event, error)
}

type Processor interface {
	Process(Event) error
}

type Type int

const (
	Unknown Type = iota
	Message
)

type Event struct {
	Type Type
}
