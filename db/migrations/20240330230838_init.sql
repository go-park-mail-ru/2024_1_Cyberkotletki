-- +goose Up

-- Создание таблицы static
CREATE TABLE IF NOT EXISTS static
(
    id         INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    name       TEXT
        CONSTRAINT upload_name_length CHECK (LENGTH(name) <= 255) NOT NULL,
    path       TEXT
        CONSTRAINT upload_path_length CHECK (LENGTH(path) <= 255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы country
CREATE TABLE IF NOT EXISTS country
(
    id   INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    name TEXT
        CONSTRAINT name_length CHECK (LENGTH(name) <= 64) NOT NULL UNIQUE
);

-- Создание таблицы genre
CREATE TABLE IF NOT EXISTS genre
(
    id   INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    name TEXT
        CONSTRAINT genre_name_length CHECK (LENGTH(name) <= 32) UNIQUE
);

-- Создание таблицы person
CREATE TABLE IF NOT EXISTS person
(
    id              INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    first_name      TEXT
        CONSTRAINT person_first_name_length CHECK (LENGTH(first_name) <= 30) NOT NULL,
    last_name       TEXT
        CONSTRAINT person_last_name_length CHECK (LENGTH(last_name) <= 30)  NOT NULL,
    birth_date      TIMESTAMP,
    death_date      TIMESTAMP CHECK (birth_date IS NULL OR death_date > birth_date),
    start_career    TIMESTAMP,
    end_career      TIMESTAMP CHECK (start_career IS NULL OR end_career > start_career),
    sex             CHAR(1)
        CONSTRAINT person_gender CHECK (sex = 'M' OR sex = 'F')             NOT NULL,
    height          INT
        CONSTRAINT height_positive CHECK (height > 50),
    spouse          TEXT
        CONSTRAINT person_spouse_length CHECK (LENGTH(spouse) <= 50),
    children        TEXT
        CONSTRAINT person_children_length CHECK (LENGTH(children) <= 150),
    photo_upload_id INT,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (photo_upload_id) REFERENCES static (id) ON DELETE SET NULL
);

-- Создание таблицы content
CREATE TABLE IF NOT EXISTS content
(
    id               INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    title            TEXT
        CONSTRAINT title_length CHECK (LENGTH(title) <= 150)               NOT NULL,
    -- original_title может повторяться, это просто название контента на языке той страны, в котором он был произведен
    original_title   TEXT
        CONSTRAINT original_title_length CHECK (LENGTH(original_title) <= 150),
    slogan           TEXT
        CONSTRAINT slogan_length CHECK (LENGTH(slogan) <= 150),
    budget           INT
        CONSTRAINT content_budget_positive CHECK (budget > 0),
    age_restriction  INT
        CONSTRAINT content_age_restriction_not_negative CHECK (age_restriction >= 0) NOT NULL,
    audience         INT
        CONSTRAINT content_audience_positive CHECK (audience >= 0),
    imdb             DECIMAL(3, 1)
        CONSTRAINT content_imdb CHECK (imdb >= 0 AND imdb <= 10) NOT NULL,
    description      TEXT
        CONSTRAINT description_length CHECK (LENGTH(description) <= 10000) NOT NULL,
    poster_upload_id INT                                                   NOT NULL,
    box_office       INT
        CONSTRAINT content_box_office_positive CHECK (box_office >= 0),
    marketing_budget INT
        CONSTRAINT content_marketing_budget_positive CHECK (marketing_budget >= 0),
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (poster_upload_id) REFERENCES static (id) ON DELETE CASCADE
);

-- Создание таблицы roles. Оно во множественном числе, тк role - ключевое слово в SQL
CREATE TABLE IF NOT EXISTS roles
(
    id   INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    name TEXT
        CONSTRAINT role_name_length CHECK (LENGTH(name) <= 30) UNIQUE NOT NULL
);


-- Создание таблицы person_roles
CREATE TABLE IF NOT EXISTS person_roles
(
    id         INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    roles_id    INT NOT NULL,
    person_id  INT NOT NULL,
    content_id INT NOT NULL,
    FOREIGN KEY (roles_id) REFERENCES roles (id) ON DELETE CASCADE,
    FOREIGN KEY (person_id) REFERENCES person (id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES content (id) ON DELETE CASCADE,
    CONSTRAINT roles_unique UNIQUE (person_id, content_id, roles_id)
);

-- Создание таблицы movie
CREATE TABLE IF NOT EXISTS movie
(
    id         INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    content_id INT NOT NULL UNIQUE,
    premiere   TIMESTAMP,
    release    TIMESTAMP,
    duration   INT CHECK (duration > 0),
    FOREIGN KEY (content_id) REFERENCES content (id) ON DELETE CASCADE
);

-- Создание таблицы series
CREATE TABLE IF NOT EXISTS series
(
    id         INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    year_start INT
        CONSTRAINT series_year_start_positive CHECK (year_start > 0) NOT NULL,
    year_end   INT CHECK (year_end >= year_start)
        CONSTRAINT series_year_end_positive CHECK (year_end > 0),
    content_id INT                                                   NOT NULL,
    FOREIGN KEY (content_id) REFERENCES content (id) ON DELETE CASCADE
);

-- Создание таблицы season
CREATE TABLE IF NOT EXISTS season
(
    id         INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    series_id  INT                                                   NOT NULL,
    year_start INT
        CONSTRAINT season_year_start_positive CHECK (year_start > 0) NOT NULL,
    year_end   INT CHECK (year_end >= year_start)
        CONSTRAINT season_year_end_positive CHECK (year_end > 0),
    FOREIGN KEY (series_id) REFERENCES series (id) ON DELETE CASCADE
);

-- Создание таблицы episode
CREATE TABLE IF NOT EXISTS episode
(
    id             INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    season_id      INT                                                        NOT NULL,
    episode_number INT
        CONSTRAINT episode_episode_number_positive CHECK (episode_number > 0) NOT NULL,
    title          TEXT
        CONSTRAINT episode_title_length CHECK (LENGTH(title) <= 150)         NOT NULL,
    FOREIGN KEY (season_id) REFERENCES season (id) ON DELETE CASCADE
);

-- Создание таблицы users. Оно во множественном числе, тк user - ключевое слово в SQL
CREATE TABLE IF NOT EXISTS users
(
    id               INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    name             TEXT
        CONSTRAINT user_name_length CHECK (LENGTH(name) <= 30),
    -- email UNIQUE, то есть может служить альтернативным ключом
    email            TEXT
        CONSTRAINT email_length CHECK (LENGTH(email) <= 256) UNIQUE                   NOT NULL,
    password_hashed  bytea
        CONSTRAINT password_hashed_length CHECK (OCTET_LENGTH(password_hashed) <= 32) NOT NULL,
    -- salt_password - поле используется для увеличения безопасности паролей(предотвращает использование одинаковых хэшей, тк примешивается к
    -- нему, понижает возможность взлома пароля путем перебора)
    -- когда пользователь входит в систему, система берет введенный пароль, добавляет к нему соль из базы данных, хеширует комбинацию и сравнивает
    -- полученный хеш с хешем, хранящимся в базе данных. Если хеши совпадают, пароль считается верным.
    salt_password    bytea
        CONSTRAINT salt_password_length CHECK (OCTET_LENGTH(salt_password) <= 8)      NOT NULL,
    avatar_upload_id INT,
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (avatar_upload_id) REFERENCES static (id) ON DELETE SET NULL
);

-- Создание таблицы review
CREATE TABLE IF NOT EXISTS review
(
    id             INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    user_id        INT                                                                       NOT NULL,
    content_id     INT                                                                       NOT NULL,
    title          TEXT
        CONSTRAINT review_title_length CHECK (LENGTH(title) > 0 AND LENGTH(title) <= 50)     NOT NULL,
    text           TEXT
        CONSTRAINT review_text_length CHECK (LENGTH(text) > 0 AND LENGTH(text) <= 10000)     NOT NULL,
    content_rating INT
        CONSTRAINT rating_user_positive CHECK (content_rating >= 1 AND content_rating <= 10) NOT NULL,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES content (id) ON DELETE CASCADE,
    CONSTRAINT user_content_unique UNIQUE (user_id, content_id)
);

-- Создание таблицы compilation_type
CREATE TABLE IF NOT EXISTS compilation_type
(
    id   INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    type TEXT
        CONSTRAINT compilation_type_length CHECK (LENGTH(type) <= 30) UNIQUE NOT NULL
);

-- Создание таблицы compilation
CREATE TABLE IF NOT EXISTS compilation
(
    id                  INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    title               TEXT
        CONSTRAINT compilation_title_length CHECK (LENGTH(title) > 0 AND LENGTH(title) <= 50) NOT NULL,
    compilation_type_id INT
        CONSTRAINT compilation_type_id_positive CHECK (compilation_type_id > 0)               NOT NULL,
    poster_upload_id    INT,
    FOREIGN KEY (poster_upload_id) REFERENCES static (id) ON DELETE SET NULL,
    FOREIGN KEY (compilation_type_id) REFERENCES compilation_type (id) ON DELETE CASCADE,
    CONSTRAINT compilation_type_unique UNIQUE (id, compilation_type_id)
);

-- Создание таблицы compilation_content
CREATE TABLE IF NOT EXISTS compilation_content
(
    compilation_id INT NOT NULL,
    content_id     INT NOT NULL,
    PRIMARY KEY (compilation_id, content_id),
    FOREIGN KEY (compilation_id) REFERENCES compilation (id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES content (id) ON DELETE CASCADE
);

-- Создание таблицы review_like
CREATE TABLE IF NOT EXISTS review_like
(
    review_id  INT     NOT NULL,
    user_id    INT     NOT NULL,
    value      BOOLEAN NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (review_id, user_id),
    FOREIGN KEY (review_id) REFERENCES review (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Создание таблицы genre_content
CREATE TABLE IF NOT EXISTS genre_content
(
    genre_id   INT NOT NULL,
    content_id INT NOT NULL,
    PRIMARY KEY (genre_id, content_id),
    FOREIGN KEY (genre_id) REFERENCES genre (id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES content (id) ON DELETE CASCADE
);

-- Создание таблицы country_content
CREATE TABLE IF NOT EXISTS country_content
(
    country_id INT NOT NULL,
    content_id INT NOT NULL,
    PRIMARY KEY (country_id, content_id),
    FOREIGN KEY (country_id) REFERENCES country (id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES content (id) ON DELETE CASCADE
);

-- Создание таблицы content_type
CREATE TABLE IF NOT EXISTS content_type
(
    id         INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    content_id INT                                                      NOT NULL,
    type       TEXT
        CONSTRAINT type_check CHECK (type = 'movie' OR type = 'series') NOT NULL,
    FOREIGN KEY (content_id) REFERENCES content (id) ON DELETE CASCADE
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +goose StatementEnd

-- Добавление триггера к каждой таблице, которая имеет поле updated_at
CREATE TRIGGER update_at_person
    BEFORE UPDATE
    ON person
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_at_content
    BEFORE UPDATE
    ON content
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_at_users
    BEFORE UPDATE
    ON users
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_at_review
    BEFORE UPDATE
    ON review
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_at_review_like
    BEFORE UPDATE
    ON review_like
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Индекс для быстрой сортировки рецензий по количеству лайков
CREATE INDEX idx_review_like_review_id ON review_like (review_id);
-- Индекс для быстрого подсчёта рейтинга контента по количеству лайков
CREATE INDEX idx_review_content_id ON review (content_id);

-- Роли по умолчанию
INSERT INTO roles (name) VALUES ('actor'), ('director'), ('producer'), ('writer'), ('operator'), ('composer'), ('editor');
