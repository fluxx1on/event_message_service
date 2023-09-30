CREATE TYPE client_tag AS ENUM ('silver', 'gold', 'vip');

CREATE TYPE filter_attr AS ENUM('tag', 'code');

CREATE TABLE IF NOT EXISTS mailing (
    id SERIAL PRIMARY KEY,
    message_text TEXT NOT NULL,
    mobile_operator_code INTEGER NOT NULL,
    tag client_tag NOT NULL,
    filter_choice filter_attr NOT NULL,
    datetime_start TIMESTAMP NOT NULL,
    datetime_end TIMESTAMP NOT NULL,
    interval_start TIMESTAMP NOT NULL,
    interval_end TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS client (
    id SERIAL PRIMARY KEY,
    phone_number BIGINT UNIQUE NOT NULL,
    mobile_operator_code INTEGER NOT NULL,
    tag client_tag NOT NULL,
    time_zone INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS message (
    id SERIAL PRIMARY KEY,
    date_time_creation TIMESTAMP NOT NULL DEFAULT now(),
    delivery_status BOOLEAN NOT NULL,
    mailing_id BIGINT REFERENCES mailing(id),
    client_id BIGINT REFERENCES client(id)
);

