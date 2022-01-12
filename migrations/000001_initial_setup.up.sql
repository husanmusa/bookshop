CREATE TABLE IF NOT EXISTS authors
(
    id         uuid        not null PRIMARY KEY,
    name       VARCHAR(64) not null,
    created_at TIMESTAMP   not null DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP   not null DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP
);

CREATE TABLE IF NOT EXISTS books
(
    id         uuid PRIMARY KEY not null,
    name       VARCHAR(64)      not null,
    author_id  uuid             not null references authors (id),
    created_at TIMESTAMP        not null DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP        not null DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP
);

CREATE TABLE IF NOT EXISTS categories
(
    id         uuid        not null PRIMARY KEY,
    name       VARCHAR(64) not null,
    parent_id  uuid,
    created_at TIMESTAMP   not null DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP   not null DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP
);

create table if not exists book_categories
(
    id          serial primary key not null,
    book_id     uuid references books (id),
    category_id uuid references categories (id)
);
