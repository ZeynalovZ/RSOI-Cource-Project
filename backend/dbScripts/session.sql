create database sessions;
CREATE TABLE credentials (
    id uuid PRIMARY KEY,
    login VARCHAR(80) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    session_id uuid NOT NULL,
    role VARCHAR(80) NOT NULL CHECK (role in ('admin', 'user', 'prime_user'))
);

CREATE TABLE sessions (
    id uuid PRIMARY KEY,
    refresh_token VARCHAR(255) NOT NULL,
    access_token VARCHAR(255) NOT NULL,
    expires_at DATE NOT NULL,
    user_id uuid NOT NULL,
    is_authenticated boolean NOT NULL
);
