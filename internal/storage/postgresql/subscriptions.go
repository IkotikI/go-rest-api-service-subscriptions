package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	microservice "github.com/ikotiki/go-rest-api-service-subscriptions"
	"github.com/rs/zerolog/log"

	"github.com/ikotiki/go-rest-api-service-subscriptions/internal/storage"
	"github.com/ikotiki/go-rest-api-service-subscriptions/pkg/e"

	"github.com/jmoiron/sqlx"
)

const (
	errStrUserSubscriptionPairAlreadyExists = "duplicate key value violates unique constraint \"subscriptions_user_id_service_name_key\""
)

type SubscriptionsStore struct {
	db      *sqlx.DB
	builder *postgresSQLBuilder
}

func NewSubscriptionsStore(store *SQLStorage) *SubscriptionsStore {
	return &SubscriptionsStore{db: store.db, builder: &postgresSQLBuilder{store.builder}}
}

func (s *SubscriptionsStore) GetByID(ctx context.Context, id microservice.SubscriptionID) (sub *microservice.Subscription, err error) {
	const op = "storage.postgresql.subscriptions.getbyid"
	sub = &microservice.Subscription{}

	q := sprintf(`SELECT * FROM %s WHERE id = $1`, TableSubscriptions)

	log.Debug().Str("query", q).Int("id", int(id)).Msg(op)

	err = s.db.GetContext(ctx, sub, q, id)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSuchSubscription
	}

	return sub, err
}

func (s *SubscriptionsStore) Create(ctx context.Context, sub *microservice.Subscription) (id microservice.SubscriptionID, err error) {
	const op = "storage.postgresql.subscriptions.create"
	q := sprintf(`
		INSERT INTO %s (user_id, service_name, monthly_price, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, TableSubscriptions)

	log.Debug().Str("query", q).Interface("subscription", sub).Msg(op)

	row := s.db.QueryRowxContext(ctx, q, sub.UserID, sub.ServiceName, sub.MonthlyPrice, sub.StartDate, sub.EndDate)
	if err = row.Err(); err != nil {
		if e.HasText(err, errStrUserSubscriptionPairAlreadyExists) {
			return 0, storage.ErrUserSubscriptionPairAlreadyExists
		}
		return 0, e.Wrap(op, err)
	}

	row.Scan(&id)
	if id == 0 {
		return 0, e.Wrap(op, err)
	}

	return id, nil
}

func (s *SubscriptionsStore) Update(ctx context.Context, sub *microservice.Subscription) (err error) {
	const op = "storage.postgresql.subscriptions.update"
	q := sprintf(`
		UPDATE %s SET (service_name, monthly_price, start_date, end_date) = ($2, $3, $4, $5)
		WHERE id = $1
	`, TableSubscriptions)

	log.Debug().Str("query", q).Interface("subscription", sub).Msg(op)

	_, err = s.db.ExecContext(ctx, q, sub.ID, sub.ServiceName, sub.MonthlyPrice, sub.StartDate, sub.EndDate)
	return e.WrapIfErr(op, err)
}

func (s *SubscriptionsStore) DeleteByID(ctx context.Context, id microservice.SubscriptionID) (err error) {
	const op = "storage.postgresql.subscriptions.deletebyid"
	q := sprintf(`
		DELETE FROM %s WHERE id = $1
	`, TableSubscriptions)

	log.Debug().Str("query", q).Int("id", int(id)).Msg(op)

	res, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return e.Wrap(op, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return e.Wrap(fmt.Sprintf("%s.rows_affected", op), err)
	}
	if n == 0 {
		return e.Wrap(op, storage.ErrNoSuchSubscription)
	}

	return nil
}

func (s *SubscriptionsStore) Query(ctx context.Context, args *storage.QueryArgs) (subs []*microservice.Subscription, err error) {
	const op = "storage.postgresql.subscriptions.query"
	q := sprintf(`SELECT * FROM %s `, TableSubscriptions)

	// Custom handling for where statement
	where, queryArgs := s.builder.buildWhere(args)
	q += where

	// For other use builder
	queryEnd, queryArgs2 := s.builder.buildParts([]string{"group_by", "order_by", "limit"}, args, 2)
	q += queryEnd
	queryArgs = append(queryArgs, queryArgs2...)

	log.Debug().Str("query", q).Interface("queryArgs", queryArgs).Msg(op)

	err = s.db.SelectContext(ctx, &subs, q, queryArgs...)
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	return subs, nil
}

func (s *SubscriptionsStore) Sum(ctx context.Context, args *storage.QueryArgs) (sum microservice.Price, err error) {
	const op = "storage.postgresql.subscriptions.sum"
	q := sprintf(`SELECT sum(monthly_price) FROM %s AS sum `, TableSubscriptions)

	where, queryArgs := s.builder.buildWhere(args)
	q += where

	log.Debug().Str("query", q).Interface("queryArgs", queryArgs).Msg(op)

	structSum := struct {
		Sum microservice.Price
	}{
		Sum: 0,
	}

	err = s.db.GetContext(ctx, &structSum, q, queryArgs...)
	if err != nil {
		if e.HasText(err, "converting NULL to int is unsupported") {
			return 0, storage.ErrNoSuchSubscription
		}
		return 0, e.Wrap(op, err)
	}
	return structSum.Sum, nil
}
