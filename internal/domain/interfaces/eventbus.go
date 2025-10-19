package interfaces

import "github.com/throindev/payments/internal/domain"

type EventBus interface {
	Publish(topic string, body *domain.Event) error
	Consume(topic string, handler func(body domain.Event)) error
}
