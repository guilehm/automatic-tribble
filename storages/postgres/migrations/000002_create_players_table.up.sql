CREATE TABLE players
(
    id         serial PRIMARY KEY,
    name       varchar(32)  NOT NULL,
    user_id    int          NOT NULL,
    xp         bigint       NOT NULL,
    sprite     varchar(128) NOT NULL,
    position_x smallint     NOT NULL,
    position_y smallint     NOT NULL
);

ALTER TABLE players
    ADD CONSTRAINT players_user_id_fk_user_id
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;

CREATE INDEX players_user_id ON players (user_id);
CREATE UNIQUE INDEX name_unique_players_idx on players (LOWER(name));
