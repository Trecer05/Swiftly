-- тут две таблицы для будущего функционала команд
-- CREATE TABLE IF NOT EXISTS commands (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(100) NOT NULL,
--     description VARCHAR(200) NOT NULL,
--     created_at TIMESTAMP DEFAULT NOW()
-- );

-- CREATE TABLE IF NOT EXISTS command_projects (
--     command_id INTEGER NOT NULL,
--     project_id INTEGER NOT NULL,
--     PRIMARY KEY (command_id, project_id),
--     FOREIGN KEY (command_id) REFERENCES commands(id) ON DELETE CASCADE,
--     FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
-- );

-- CREATE TABLE IF NOT EXISTS command_users (
--     command_id INTEGER NOT NULL,
--     user_id INTEGER NOT NULL,
--     PRIMARY KEY (command_id, user_id),
--     FOREIGN KEY (command_id) REFERENCES commands(id) ON DELETE CASCADE,
--     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
-- );

CREATE TABLE IF NOT EXISTS users (
    id integer PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    username VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(200) NOT NULL,
    avatar_url TEXT,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);

CREATE TABLE IF NOT EXISTS chats (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS chat_users (
    chat_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    PRIMARY KEY (chat_id, user_id),
    FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(200) NOT NULL,
    owner_id BIGINT NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS group_users (
    group_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    PRIMARY KEY (group_id, user_id),
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS chat_messages (
    id SERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    text TEXT NOT NULL,
    read boolean DEFAULT false,
    edited boolean DEFAULT false,
    sent_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS chat_file_urls (
    id SERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    file_url TEXT,
    FOREIGN KEY (chat_id) REFERENCES chat_messages(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS chat_messages_file_urls (
    id SERIAL PRIMARY KEY,
    chat_message_id BIGINT NOT NULL,
    chat_file_id BIGINT NOT NULL,
    FOREIGN KEY (chat_file_id) REFERENCES chat_file_urls(id) ON DELETE CASCADE,
    FOREIGN KEY (chat_message_id) REFERENCES chat_messages(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS group_messages (
    id SERIAL PRIMARY KEY,
    group_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    text TEXT NOT NULL,
    read boolean DEFAULT false,
    edited boolean DEFAULT false,
    sent_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS group_file_urls (
    id SERIAL PRIMARY KEY,
    group_id BIGINT NOT NULL,
    file_url TEXT,
    FOREIGN KEY (group_id) REFERENCES group_messages(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS group_messages_file_urls (
    id SERIAL PRIMARY KEY,
    group_message_id BIGINT NOT NULL,
    group_file_id BIGINT NOT NULL,
    FOREIGN KEY (group_file_id) REFERENCES group_file_urls(id) ON DELETE CASCADE,
    FOREIGN KEY (group_message_id) REFERENCES group_messages(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(200) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users_projects (
    project_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'участник',
    is_admin BOOLEAN DEFAULT FALSE,
    added_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (project_id, user_id),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TYPE ApplicationStatus AS ENUM ('pending', 'accepted', 'rejected');

CREATE TABLE IF NOT EXISTS team_applications (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    project_id BIGINT NOT NULL,
    status ApplicationStatus NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS project_invites (
    id SERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL,
    invite_code VARCHAR(100) UNIQUE NOT NULL,
    creator_id BIGINT NOT NULL,
    expires_at TIMESTAMP NOT NULL DEFAULT NOW() + INTERVAL '3 days',
    created_at TIMESTAMP DEFAULT NOW(),
    is_single_use BOOLEAN NOT NULL DEFAULT false,
    used BOOLEAN NOT NULL DEFAULT false,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX invite_code_idx ON project_invites(invite_code);

CREATE INDEX idx_users_projects_project_user ON users_projects(project_id, user_id);
CREATE INDEX idx_users_projects_user_project ON users_projects(user_id, project_id);