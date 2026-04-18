-- +goose Up
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS requests (
    id SERIAL PRIMARY KEY,
    status BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS files (
    id SERIAL PRIMARY KEY,
    request_id FOREIGN KEY REFERENCES requests(id) ON DELETE CASCADE,
    filepath TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS pdn (
    id SERIAL PRIMARY KEY,
    file_id FOREIGN KEY REFERENCES files(id) ON DELETE CASCADE,
    type_of_pdn TEXT,
    start_time REAL,
    end_time REAL
);

CREATE TABLE IF NOT EXISTS transcriptions (
    id SERIAL PRIMARY KEY,
    request_id FOREIGN KEY REFERENCES requests(id) ON DELETE CASCADE,
    original_text TEXT,
    anon_text TEXT
);

-- +goose Down
SELECT 'down SQL query';

DROP TABLE IF EXISTS pdn;
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS transcriptions;
DROP TABLE IF EXISTS requests;