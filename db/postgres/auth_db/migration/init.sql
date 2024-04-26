create table IF NOT EXISTS sessions
(
    id uuid not null primary key,
    user_id uuid not null,
    access_token text not null,
    refresh_token text not null,
    access_token_expired_at timestamptz not null,
    refresh_token_expired_at timestamptz not null,
    created_at timestamptz not null,
    updated_at timestamptz not null,
    deleted_at timestamptz
);
