CREATE TABLE favorite_phrases (
    id INTEGER PRIMARY KEY,
    category_id INTEGER NOT NULL,
    content TEXT NOT NULL CHECK (length(trim(content)) > 0),
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
    UNIQUE (category_id, content)
);

CREATE TABLE generation_history (
    id INTEGER PRIMARY KEY,
    category_id INTEGER NOT NULL,
    content TEXT NOT NULL CHECK (length(trim(content)) > 0),
    generated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

CREATE INDEX idx_favorite_phrases_category_id ON favorite_phrases(category_id);
CREATE INDEX idx_generation_history_category_id ON generation_history(category_id);
CREATE INDEX idx_generation_history_generated_at ON generation_history(generated_at DESC);
