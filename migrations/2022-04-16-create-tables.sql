CREATE TABLE IF NOT EXISTS users
(
    login         TEXT PRIMARY KEY,
    password_hash TEXT        NOT NULL,
    session       TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS orders
(
    id          SERIAL PRIMARY KEY,
    order_num   TEXT UNIQUE              NOT NULL,
    login       TEXT                     NOT NULL REFERENCES users (login),
    uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL,
    status      TEXT                     NOT NULL,
    accrual     NUMERIC
);

CREATE TABLE IF NOT EXISTS balance
(
    id      SERIAL PRIMARY KEY,
    login   TEXT     NOT NULL REFERENCES users (login),
    balance NUMERIC  NOT NULL DEFAULT 0,
    spent   NUMERIC  NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS balance_withdrawals
(
    id           SERIAL PRIMARY KEY,
    login        TEXT                     NOT NULL REFERENCES users (login),
    order_num    TEXT                     NOT NULL,
    sum          NUMERIC                  NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE NOT NULL
);