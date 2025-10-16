package domain

type SubscriptionRepository interface {
	Create(subscription *Subscription) error
	Update(subscription *Subscription) error
	FindByID(id string) (*Subscription, error)
	FindByUserID(userID string) (*Subscription, error)
}
