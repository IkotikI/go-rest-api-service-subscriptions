-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE subscriptions (
    id serial PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL,
    service_name varchar(120) NOT NULL,
    monthly_price integer CHECK (monthly_price >= 0) NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP,
    UNIQUE (user_id, service_name)
);

CREATE INDEX idx_subscriptions_user_id_service_name ON subscriptions(user_id, service_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_subscriptions_user_id_service_name;
DROP TABLE IF EXISTS subscriptions;
DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd
