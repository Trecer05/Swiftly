CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TYPE file_visibility AS ENUM ('public', 'private', 'shared');

CREATE TABLE folders (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    owner_id INTEGER,  -- Владелец (пользователь ИЛИ команда)
    owner_type VARCHAR(10) CHECK (owner_type IN ('user', 'team')),
    parent_folder_id UUID REFERENCES folders(uuid) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT valid_owner CHECK (
        (owner_type = 'user' AND owner_id IS NOT NULL) OR
        (owner_type = 'team' AND owner_id IS NOT NULL)
    )
);

CREATE TABLE files (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    folder_id UUID REFERENCES folders(uuid) ON DELETE CASCADE,
    storage_path VARCHAR(500) NOT NULL,
    original_filename VARCHAR(500) NOT NULL,
    display_name VARCHAR(500),
    mime_type VARCHAR(255),
    size BIGINT NOT NULL DEFAULT 0,
    visibility file_visibility NOT NULL DEFAULT 'private',
    
    created_by INTEGER NOT NULL, -- Кто загрузил файл
    owner_id INTEGER, -- Владелец (пользователь ИЛИ команда)
    owner_type VARCHAR(10) CHECK (owner_type IN ('user', 'team')),
    
    uploaded_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    hash VARCHAR(64),
    version INTEGER DEFAULT 1,
    
    CONSTRAINT valid_file_owner CHECK (
        (owner_type = 'user' AND owner_id IS NOT NULL) OR
        (owner_type = 'team' AND owner_id IS NOT NULL)
    )
);

CREATE TABLE shared_access (
    id SERIAL PRIMARY KEY,
    file_id UUID REFERENCES files(uuid) ON DELETE CASCADE,
    folder_id UUID REFERENCES folders(uuid) ON DELETE CASCADE,
    
    shared_with_id INTEGER NOT NULL,
    shared_with_type VARCHAR(10) CHECK (shared_with_type IN ('user', 'team')),
    
    shared_by INTEGER NOT NULL,
    shared_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP,
    
    CONSTRAINT shared_target_check CHECK (
        (file_id IS NOT NULL AND folder_id IS NULL) OR
        (file_id IS NULL AND folder_id IS NOT NULL)
    )
);

CREATE TABLE file_tags (
    file_id UUID REFERENCES files(uuid) ON DELETE CASCADE,
    tag VARCHAR(100) NOT NULL,
    added_by INTEGER NOT NULL,
    added_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (file_id, tag)
);

CREATE TABLE file_versions (
    id SERIAL PRIMARY KEY,
    file_id UUID REFERENCES files(uuid) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    storage_path VARCHAR(500) NOT NULL,
    size BIGINT NOT NULL,
    hash VARCHAR(64),
    created_by INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    comment TEXT,
    UNIQUE(file_id, version)
);

CREATE INDEX idx_folders_owner ON folders(owner_type, owner_id);
CREATE INDEX idx_shared_access ON shared_access(shared_with_type, shared_with_id);

CREATE INDEX idx_files_folder ON files(folder_id);
CREATE INDEX idx_files_owner ON files(owner_type, owner_id);
CREATE INDEX idx_files_visibility ON files(visibility);
CREATE INDEX idx_files_display_name ON files USING gin(display_name gin_trgm_ops); -- Для поиска
CREATE INDEX idx_files_uploaded ON files(uploaded_at DESC);

CREATE INDEX idx_folders_parent ON folders(parent_folder_id);
CREATE INDEX idx_folders_name ON folders USING gin(name gin_trgm_ops);

CREATE INDEX idx_file_tags_tag ON file_tags(tag);
CREATE INDEX idx_file_tags_file ON file_tags(file_id);

CREATE INDEX idx_files_user_folder ON files(owner_id, folder_id) WHERE owner_type = 'user';

CREATE INDEX idx_files_team_visibility ON files(owner_id, visibility) WHERE owner_type = 'team';