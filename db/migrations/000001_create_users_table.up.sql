CREATE TABLE users
(
    id            serial PRIMARY KEY,
    name          varchar(50)              NOT NULL,
    email         varchar(512)             NOT NULL,
    password      varchar(256)             NOT NULL,
    token         varchar(256)             NOT NULL,
    refresh_token varchar(256)             NOT NULL,
    date_joined   timestamp with time zone NOT NULL
);

CREATE UNIQUE INDEX name_unique_users_idx on users (LOWER(name));
CREATE UNIQUE INDEX email_unique_users_idx on users (LOWER(email));
