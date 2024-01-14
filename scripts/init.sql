CREATE TABLE IF NOT EXISTS users
(
    id   VARCHAR(255) PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS wallets
(
    id       VARCHAR(255) PRIMARY KEY,
    user_id  VARCHAR(255) REFERENCES users (id) ON DELETE CASCADE,
    balance  DECIMAL(10, 2) DEFAULT 0.00,
    currency VARCHAR(3)     DEFAULT 'EUR',
    wallet_version  INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS transactions
(
    id               VARCHAR(255) PRIMARY KEY,
    wallet_id        VARCHAR(255),
    user_id          VARCHAR(255),
    amount           DECIMAL(10, 2),
    transaction_type VARCHAR(255),
    status           VARCHAR(255) DEFAULT 'PENDING',
    created_at       VARCHAR(255),
    updated_at       VARCHAR(255)
);

-- //TODO: CREATE INDEX

-- Insert sample users
INSERT INTO users (id, name)
VALUES ('1', 'Jane'),
       ('2', 'John'),
       ('3', 'Peter');

-- INSERT INTO wallets (user_id, balance, currency) VALUES
--                                                      (1, 100.00, 'USD', 0),
--                                                      (2, 50.00, 'EUR', 0),
--                                                      (3, 200.00, 'GBP', 0);
