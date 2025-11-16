CREATE TYPE IF NOT EXISTS priority_level AS ENUM ('low', 'medium', 'high');

CREATE TABLE priority_definitions (
    level priority_level PRIMARY KEY,
    title TEXT NOT NULL,
    color VARCHAR(7) NOT NULL CHECK (color ~ '^#[0-9A-Fa-f]{6}$')
);

INSERT INTO priority_definitions VALUES 
('low', 'Low', '#00FF00'),
('medium', 'Medium', '#FFFF00'),
('high', 'High', '#FF0000');

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    author_id INTEGER NOT NULL,
    developer_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    label VARCHAR(100),
    priority priority_level REFERENCES priority_definitions(level),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    completed_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE project_tasks (
    project_id INTEGER NOT NULL,
    task_id INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'todo',
    FOREIGN KEY (task_id) REFERENCES tasks(id)
);
