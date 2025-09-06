CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE subscriptions (
    id serial PRIMARY KEY AUTOINCREMENT,
    user_id UUID DEFAULT uuid_generate_v4(),
    service_name varchar(120) NOT NULL,
    monthly_price integer NOT NULL CHECK (monthly_price >= 0)
    start_date DATE NOT NULL,
    end_date DATE,
)

CREATE INDEX idx_subscriptions_id ON subscriptions(id);
CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);

DROP TABLE IF EXISTS subscriptions;