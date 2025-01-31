CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    CHECK (id > 0)
);

CREATE TABLE IF NOT EXISTS products (
    id INTEGER PRIMARY KEY,
    CHECK (id > 0)
);

CREATE TABLE IF NOT EXISTS keyWords (
    id SERIAL PRIMARY KEY,
    kw_name VARCHAR(255) NOT NULL,
    CHECK (LENGTH(TRIM(kw_name)) > 0)
);

CREATE TABLE IF NOT EXISTS user_kw (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    kw_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (kw_id) REFERENCES keyWords(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS product_kw (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL,
    kw_id INTEGER NOT NULL,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (kw_id) REFERENCES keyWords(id) ON DELETE CASCADE
);