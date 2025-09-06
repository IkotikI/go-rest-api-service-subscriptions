package storage

import (
	"errors"

	"github.com/google/uuid"
)

var ErrNoSuchSubscription = errors.New("no such subscription")
var ErrSubscriptionAlreadyExists = errors.New("subscription already exists")
var ErrNoUserID = errors.New("no user is provided or its invalid")
var ErrNoSubscriptionID = errors.New("no subscription is provided or its invalid")
var ErrUserSubscriptionPairAlreadyExists = errors.New("user-subscription pair already exists")

type UserID = uuid.UUID
type SubscriptionID = int64
type Price int

/* ---- Subscription Type ---- */
type Subscription struct {
	ID           SubscriptionID `json:"id" db:"id" swaggerignore:"true"`
	UserID       UserID         `json:"user_id" db:"user_id" binding:"required"`
	ServiceName  string         `json:"service_name" db:"service_name" binding:"required"`
	MonthlyPrice Price          `json:"monthly_price" db:"monthly_price" binding:"required"`
	StartDate    Date           `json:"start_date" db:"start_date" binding:"required"`
	EndDate      Date           `json:"end_date,omitempty,omitzero" db:"end_date"`
}

/* ---- Query ---- */
// Provide abstract arguments for making SQL queries.
// Concrete implementation lies on chosen

/* ---- Tables ---- */
type QueryArgs struct {
	From   From          `json:"from"`
	Where  []Where       `json:"where"`
	Order  []OrderStruct `json:"order"`
	Limit  int64         `json:"limit"`
	Offset int64         `json:"offset"`
}

type From string

const (
	FromSubscriptions From = "subscriptions"
)

type Operator string

const (
	OpEqual       Operator = "="
	OpNotEqual    Operator = "!="
	OpLess        Operator = "<"
	OpMore        Operator = ">"
	OpLessOrEqual Operator = "<="
	OpMoreOrEqual Operator = ">="
	OpIn          Operator = "IN"
)

type Where struct {
	Column   string
	Operator Operator
	Value    interface{}
}

type OrderStruct struct {
	OrderBy string
	Order   Order
}

type Order string

const (
	OrderASC  Order = "ASC"
	OrderDECS Order = "DESC"
)
