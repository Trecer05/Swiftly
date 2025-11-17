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

CREATE TYPE IF NOT EXISTS task_status AS ENUM ('todo', 'in_progress', 'completed');

CREATE TABLE project_columns (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL,
    title VARCHAR(100) NOT NULL,
    position INTEGER NOT NULL,
    color VARCHAR(7) DEFAULT '#D3D3D3',
    wip_limit INTEGER,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER NOT NULL,
    
    UNIQUE(project_id, position),
    CHECK (position >= 0)
);

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    author_id INTEGER NOT NULL,
    developer_id INTEGER NOT NULL,
    column_id INTEGER NOT NULL REFERENCES project_columns(id),
    position_in_column INTEGER NOT NULL DEFAULT 0;
    title TEXT NOT NULL,
    description TEXT,
    label VARCHAR(100),
    status task_status NOT NULL DEFAULT 'todo',
    priority priority_level REFERENCES priority_definitions(level),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    completed_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE project_tasks (
    project_id INTEGER NOT NULL,
    task_id INTEGER NOT NULL,
    FOREIGN KEY (task_id) REFERENCES tasks(id)
);

-- Для быстрого получения задач колонки в правильном порядке
CREATE INDEX idx_tasks_column_position ON tasks(column_id, position_in_column);

-- Для поиска задач по проекту (через колонку)
CREATE INDEX idx_project_columns_project ON project_columns(project_id);

-- Для статистики и отчетов
CREATE INDEX idx_tasks_priority ON tasks(priority);
CREATE INDEX idx_tasks_created_at ON tasks(created_at);
CREATE INDEX idx_tasks_completed_at ON tasks(completed_at);
