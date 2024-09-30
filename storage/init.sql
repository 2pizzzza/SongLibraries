CREATE TABLE IF NOT EXISTS songs
(
    id           SERIAL PRIMARY KEY,
    group_name   VARCHAR(255),
    song_title   VARCHAR(255),
    release_date DATE DEFAULT CURRENT_DATE,
    lyrics       TEXT DEFAULT '',
    link         VARCHAR(255) DEFAULT ''
);
