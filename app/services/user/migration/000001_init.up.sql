CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    usr_description VARCHAR(255) NOT NULL,
    is_email_verified BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS codes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    email_code INTEGER NOT NULL UNIQUE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE 
);

CREATE TABLE IF NOT EXISTS user_interests (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    interest_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (interest_id) REFERENCES interests(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS interests (
    id SERIAL PRIMARY KEY,
    interest VARCHAR(255)
);