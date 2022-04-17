CREATE TABLE users
(
    id          serial PRIMARY KEY,
    name        varchar(50)              NOT NULL,
    email       varchar                  NOT NULL,
    date_joined timestamp with time zone NOT NULL
);

CREATE UNIQUE INDEX email_unique_idx on users (LOWER(email));
