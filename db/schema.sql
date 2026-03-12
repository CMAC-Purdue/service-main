CREATE TABLE IF NOT EXISTS officers (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    title TEXT NOT NULL,
    linkedin TEXT,
    image_uri TEXT
);
