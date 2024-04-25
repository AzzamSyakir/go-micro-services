create table IF NOT EXISTS sessions
(
    id uuid not null primary key,
    user_id uuid not null,
    acces_token text not null,
    refresh_token text not null,
    total_return int not null,
    receipt_code varchar not null  unique,
   access_token_expired_at timestamptz not null,
    refresh_token_expired_at timestamptz not null,
    created_at timestamptz not null,
    updated_at timestamptz not null,
    deleted_at timestamptz
);
