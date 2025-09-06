package storage

import (
	"context"
)

type Subscriptions interface {
	Create(ctx context.Context, sub *Subscription) (id SubscriptionID, err error)
	GetByID(ctx context.Context, id SubscriptionID) (sub *Subscription, err error)
	Update(ctx context.Context, sub *Subscription) (err error)
	DeleteByID(ctx context.Context, id SubscriptionID) (err error)

	Query(ctx context.Context, args *QueryArgs) (subs []*Subscription, err error)
	Sum(ctx context.Context, args *QueryArgs) (sum Price, err error)
}

type Storage struct {
	Subscriptions
}
