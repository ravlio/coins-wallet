CREATE TYPE balance_type AS ENUM ('debit', 'credit');
DROP TYPE balance_type;
create table balances
(
    id         serial PRIMARY KEY,
    account_id int4,
    currency   int2,
    amount     int8,
    type       balance_type,
    created_at timestamp
);

drop index idx_balance_account_id_currency;

CREATE INDEX idx_balance_account_id_currency
    ON balances (account_id, currency);