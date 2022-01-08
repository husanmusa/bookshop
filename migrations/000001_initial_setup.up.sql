CREATE TABLE IF NOT EXISTS authors (
    id uuid PRIMARY KEY,
    name VARCHAR(64),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP,
    delete_at TIMESTAMP
);

