CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(30) NOT NULL,
    password VARCHAR(256) NOT NULL,
    Name VARCHAR(256) NOT NULL,
    email VARCHAR(256) NOT NULL,
    created_at BIGINT NOT NULL
);

CREATE INDEX IF NOT EXISTS login_index ON users(login);
CREATE INDEX IF NOT EXISTS email_index ON users(email);
