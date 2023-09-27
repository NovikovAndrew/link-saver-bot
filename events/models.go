package events

type Type int

const (
	Message Type = iota
	Unknown
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}
