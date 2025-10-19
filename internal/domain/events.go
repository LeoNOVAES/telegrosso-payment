package domain

import "time"

const (
	EventPaymentCreated      = "payment.created"
	EventPaymentConfirmed    = "payment.confirmed"
	EventSubscriptionExpired = "subscription.expired"
	EventUserAddedToGroup    = "user.added_to_group"
)

type Event struct {
	EventType string      `json:"event_type"`
	Data      interface{} `json:"data"`
	ProduceAt time.Time   `json:"produce_at"`
}

func NewEvent(event string, data interface{}) *Event {
	return &Event{
		EventType: event,
		Data:      data,
		ProduceAt: time.Now(),
	}
}
