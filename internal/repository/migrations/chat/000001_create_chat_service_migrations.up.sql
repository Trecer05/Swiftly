CREATE TABLE IF NOT EXISTS users (
    id integer PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    username VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(200) NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS chats (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS chat_users (
    chat_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    PRIMARY KEY (chat_id, user_id),
    FOREIGN KEY (chat_id) REFERENCES chats(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(200) NOT NULL
);

CREATE TABLE IF NOT EXISTS group_users (
    group_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    PRIMARY KEY (group_id, user_id),
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS chat_messages (
    id SERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    text TEXT NOT NULL,
    sent_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (chat_id) REFERENCES chats(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS group_messages (
    id SERIAL PRIMARY KEY,
    group_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    text TEXT NOT NULL,
    sent_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);