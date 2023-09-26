package events

type Fetcher interface {
	Fetch(limit int) Event
}

type Processor interface {
	Process(event Event)
}
