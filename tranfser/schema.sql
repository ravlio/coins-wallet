CREATE TYPE transfer_status AS ENUM ('processing','success', 'failure');

create table transfers
(
    id         serial PRIMARY KEY,
    from_id    int4,
    to_id      int4,
    currency   int2,
    amount     int8,
    status     transfer_status,
    created_at timestamp
);

CREATE INDEX idx_transfers_from_id_to_id
    ON transfers (from_id, to_id);