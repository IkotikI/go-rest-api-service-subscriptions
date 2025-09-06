package service

import (
	"context"
	"strings"
	"time"

	microservice "github.com/ikotiki/go-rest-api-service-subscriptions"
	"github.com/ikotiki/go-rest-api-service-subscriptions/internal/storage"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type SubscriptionService struct {
	store storage.Subscriptions
}

type SubscriptionQueryArgs struct {
	UserID      string  `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	ServiceName string  `json:"service_name" example:"Yandex Taxi"`
	StartDate   string  `json:"start_date" example:"2006-01-02"`
	EndDate     string  `json:"end_date" example:"2006-01-02"`
	Order       []Order `json:"order"`
}

type Order struct {
	OrderBy string `json:"order_by" example:"user_id"`
	Order   string `json:"order" example:"ASC"`
}

func NewSubscriptionService(store storage.Subscriptions) *SubscriptionService {
	return &SubscriptionService{store: store}
}

func (s *SubscriptionService) GetByID(ctx context.Context, id microservice.SubscriptionID) (sub *microservice.Subscription, err error) {
	return s.store.GetByID(ctx, id)
}

func (s *SubscriptionService) Create(ctx context.Context, sub *microservice.Subscription) (id microservice.SubscriptionID, err error) {
	return s.store.Create(ctx, sub)
}

func (s *SubscriptionService) Update(ctx context.Context, sub *microservice.Subscription) (err error) {
	return s.store.Update(ctx, sub)
}

func (s *SubscriptionService) DeleteByID(ctx context.Context, id microservice.SubscriptionID) (err error) {
	return s.store.DeleteByID(ctx, id)
}

func (s *SubscriptionService) Query(ctx context.Context, args *SubscriptionQueryArgs) (subs []*microservice.Subscription, err error) {
	queryArgs, err := s.parseQueryArgs(args)
	if err != nil {
		return nil, err
	}

	return s.store.Query(ctx, queryArgs)
}

func (s *SubscriptionService) Sum(ctx context.Context, args *SubscriptionQueryArgs) (sum microservice.Price, err error) {
	queryArgs, err := s.parseQueryArgs(args)
	if err != nil {
		return 0, err
	}

	log.Debug().Interface("queryArgs", queryArgs).Msg("query args to summation")

	return s.store.Sum(ctx, queryArgs)
}

func (s *SubscriptionService) parseQueryArgs(args *SubscriptionQueryArgs) (*storage.QueryArgs, error) {
	var queryArgs storage.QueryArgs

	if args.UserID != "" {
		userID, err := uuid.Parse(args.UserID)
		if err != nil || userID == uuid.Nil {
			return nil, ErrNoUserID
		}
		queryArgs.Where = append(queryArgs.Where, storage.Where{
			Column:   "user_id",
			Operator: storage.OpEqual,
			Value:    userID.String(),
		})
	}

	// Date start
	if args.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", args.StartDate)
		if err != nil {
			return nil, err
		}
		queryArgs.Where = append(queryArgs.Where, storage.Where{
			Column:   "start_date",
			Operator: storage.OpMoreOrEqual,
			Value:    startDate,
		})
	}

	// Date end
	if args.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", args.EndDate)
		if err != nil {
			return nil, err
		}
		queryArgs.Where = append(queryArgs.Where, storage.Where{
			Column:   "end_date",
			Operator: storage.OpLessOrEqual,
			Value:    endDate,
		})
	}

	// Order By
	orders := make([]storage.OrderStruct, 0, len(args.Order))
	for _, o := range args.Order {
		switch o.OrderBy {
		case "user_id", "service_name", "start_date", "end_date":
			queryArgsOrder := storage.OrderStruct{
				OrderBy: o.OrderBy,
			}
			// Order
			order := strings.ToUpper(o.Order)
			if order == "DESC" {
				queryArgsOrder.Order = storage.OrderDECS
			} else {
				queryArgsOrder.Order = storage.OrderASC
			}
			orders = append(orders, queryArgsOrder)
		}
	}
	queryArgs.Order = orders

	return &queryArgs, nil
}
