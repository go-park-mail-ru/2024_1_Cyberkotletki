-- +goose Up

CREATE TABLE IF NOT EXISTS favourite (
    user_id INT NOT NULL,
    content_id INT NOT NULL,
    category TEXT NOT NULL
        CHECK (category IN ('favourite', 'watching', 'watched', 'planned', 'rewatching', 'abandoned')),
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES content (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, content_id),
    UNIQUE (user_id, content_id)
);
