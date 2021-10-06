package trigger

type Subject interface {
	Attach(Observer)
	Detach(Observer)
	Update()
	Notify(event Event)
}

type Observer interface {
	OnChange(Event) bool
}

type Event interface {
	ToString() string
	GetType() EventType
}

const (
	UserLevelEvent EventType = iota + 1
)

type EventType uint32
