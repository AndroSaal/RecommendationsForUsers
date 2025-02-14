CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY, 
    CHECK (id > 0)
);

CREATE TABLE IF NOT EXISTS products (
    id INTEGER PRIMARY KEY,
    CHECK (id > 0)
);

CREATE TABLE IF NOT EXISTS user_updates (
    id SERIAL PRIMARY KEY,
    user_id INTEGER, 
    timestamp_column TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    interests VARCHAR(1024),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS product_updates (
    id SERIAL PRIMARY KEY,
    product_id INTEGER,
    timestamp_column TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    keywords VARCHAR(1024),
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);