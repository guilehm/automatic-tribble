CREATE TABLE users (
    id serial PRIMARY KEY,
    name varchar(50) NOT NULL UNIQUE,
    date_joined timestamp with time zone NOT NULL
)
