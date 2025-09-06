package microservice

import (
	"time"

	"github.com/ikotiki/go-rest-api-service-subscriptions/internal/storage"
)

type UserID = storage.UserID
type SubscriptionID = storage.SubscriptionID
type Price = storage.Price

/* ---- Subscription Type ---- */
type Subscription = storage.Subscription

type QueryArgs = storage.QueryArgs

type Date = storage.Date

func NewDate(t time.Time) Date { return storage.NewDate(t) }
