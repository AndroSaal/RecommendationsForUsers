CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    category VARCHAR(255) NOT NULL,
    prd_description VARCHAR(255) NOT NULL,
    prd_status VARCHAR(255) NOT NULL,
    CHECK (LENGTH(TRIM(category)) > 0)
);

CREATE TABLE IF NOT EXISTS keyWords (
    id SERIAL PRIMARY KEY,
    kw_name VARCHAR(255) NOT NULL, 
    CHECK (LENGTH(TRIM(kw_name)) > 0)
);

CREATE TABLE IF NOT EXISTS product_keyWord (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL,
    kw_id INTEGER NOT NULL,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (kw_id) REFERENCES keyWords(id) ON DELETE CASCADE,
    CHECK (product_id > 0)
);