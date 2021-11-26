package events

type Event_t int32
const (
    DefaultEvent Event_t = iota
)

type Event struct {
    EventType  Event_t `json:"EventType"`
    Data       string `json:"Data"`
}
