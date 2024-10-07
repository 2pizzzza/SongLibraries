CREATE TABLE IF NOT EXISTS groups
(
    id         SERIAL PRIMARY KEY,
    group_name VARCHAR(255) UNIQUE
);


CREATE TABLE IF NOT EXISTS songs
(
    id           SERIAL PRIMARY KEY,
    group_id     INT REFERENCES groups(id),
    song_title   VARCHAR(255),
    release_date VARCHAR(255),
    lyrics       TEXT DEFAULT '',
    link         VARCHAR(255) DEFAULT ''
);

