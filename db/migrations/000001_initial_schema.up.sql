-- Enable this for every SQLite connection in the application:
-- PRAGMA foreign_keys = ON;

BEGIN;

CREATE TABLE categories (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE CHECK (length(trim(name)) > 0),
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE phrase_templates (
    id INTEGER PRIMARY KEY,
    category_id INTEGER NOT NULL,
    content TEXT NOT NULL CHECK (length(trim(content)) > 0),
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

CREATE TABLE phrase_parts (
    id INTEGER PRIMARY KEY,
    category_id INTEGER NOT NULL,
    kind TEXT NOT NULL CHECK (
        kind IN ('introduction', 'subject', 'verb', 'complement', 'conclusion')
    ),
    content TEXT NOT NULL CHECK (length(trim(content)) > 0),
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

CREATE INDEX idx_phrase_templates_category_id ON phrase_templates(category_id);
CREATE INDEX idx_phrase_parts_category_id ON phrase_parts(category_id);
CREATE INDEX idx_phrase_parts_category_kind ON phrase_parts(category_id, kind);

COMMIT;
