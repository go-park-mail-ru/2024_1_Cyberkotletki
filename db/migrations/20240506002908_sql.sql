-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd


-- Создание таблицы ongoing_content
CREATE TABLE IF NOT EXISTS ongoing_content
(
    id          INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    kinopoisk_id       INT UNIQUE,
    title      TEXT
        CONSTRAINT ongoing_content_title_length CHECK (LENGTH(title) <= 150) NOT NULL,
    poster_upload_id INT,
    release_date TIMESTAMPTZ NOT NULL
    
);

CREATE TABLE IF NOT EXISTS genre_ongoing_content
(
    genre_id   INT NOT NULL,
    ongoing_content_id INT NOT NULL,
    FOREIGN KEY (genre_id) REFERENCES genre (id) ON DELETE CASCADE,
    FOREIGN KEY (ongoing_content_id) REFERENCES ongoing_content (id) ON DELETE CASCADE,
    CONSTRAINT genre_content_unique UNIQUE (genre_id, ongoing_content_id)
);

-- Создание таблицы movie
CREATE TABLE IF NOT EXISTS ongoing_movie
(
    id         INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    ongoing_content_id INT NOT NULL UNIQUE,
    premiere   TIMESTAMPTZ,
    duration   INT CHECK (duration > 0),
    FOREIGN KEY (ongoing_content_id) REFERENCES ongoing_content (id) ON DELETE CASCADE
);

-- Создание таблицы series
CREATE TABLE IF NOT EXISTS ongoing_series
(
    id         INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    year_start INT
        CONSTRAINT series_year_start_positive CHECK (year_start > 0) NOT NULL,
    year_end   INT CHECK (year_end >= year_start)
        CONSTRAINT series_year_end_positive CHECK (year_end > 0),
    ongoing_content_id INT                                                   NOT NULL,
    FOREIGN KEY (ongoing_content_id) REFERENCES ongoing_content (id) ON DELETE CASCADE
);
