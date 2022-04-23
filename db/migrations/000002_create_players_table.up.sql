CREATE TABLE players
(
    id      serial PRIMARY KEY,
    xp      bigint       NOT NULL,
    user_id int          NOT NULL,
    sprite  varchar(128) NOT NULL
);
