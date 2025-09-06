package service

import (
	"context"

	microservice "github.com/ikotiki/go-rest-api-service-subscriptions"
	"github.com/ikotiki/go-rest-api-service-subscriptions/internal/storage"
)

var (
	ErrNoSuchSubscription                = storage.ErrNoSuchSubscription
	ErrUserSubscriptionPairAlreadyExists = storage.ErrUserSubscriptionPairAlreadyExists
	ErrNoUserID                          = storage.ErrNoUserID
	ErrNoSubscriptionID                  = storage.ErrNoSubscriptionID
)

type Subscriptions interface {
	Create(ctx context.Context, sub *microservice.Subscription) (id microservice.SubscriptionID, err error)
	GetByID(ctx context.Context, id microservice.SubscriptionID) (sub *microservice.Subscription, err error)
	Update(ctx context.Context, sub *microservice.Subscription) (err error)
	DeleteByID(ctx context.Context, id microservice.SubscriptionID) (err error)

	Query(ctx context.Context, args *SubscriptionQueryArgs) (subs []*microservice.Subscription, err error)
	Sum(ctx context.Context, args *SubscriptionQueryArgs) (sum microservice.Price, err error)
}

type Service struct {
	Subscriptions
}

func NewService(store storage.Subscriptions) *Service {
	return &Service{
		Subscriptions: NewSubscriptionService(store),
	}
}
