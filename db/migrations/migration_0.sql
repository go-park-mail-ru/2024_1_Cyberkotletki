-- Создание таблицы audience
CREATE TABLE IF NOT EXISTS audience (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    content_id INT NOT NULL,
    size_in_thousands DECIMAL(10,2) CONSTRAINT size_in_thousands_positive CHECK (size_in_thousands > 0) NOT NULL,
    country_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE,
    FOREIGN KEY (country_id) REFERENCES country(id) ON DELETE CASCADE
    );

-- Создание таблицы boxoffice
CREATE TABLE IF NOT EXISTS boxoffice (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    content_id INT NOT NULL,
    country_id INT NOT NULL,
    revenue DECIMAL(10,2) CONSTRAINT revenue_positive CHECK (revenue > 0) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE,
    FOREIGN KEY (country_id) REFERENCES country(id) ON DELETE CASCADE
    );

-- Создание таблицы country
CREATE TABLE IF NOT EXISTS country (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    name TEXT CONSTRAINT name_length CHECK (LENGTH(name) <= 30) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Создание таблицы genre
CREATE TABLE IF NOT EXISTS genre (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    name TEXT CONSTRAINT genre_name_length CHECK (LENGTH(name) <= 32) UNIQUE,
    name_ru TEXT CONSTRAINT genre_name_ru_length CHECK (LENGTH(name_ru) <= 32) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы birthplace
CREATE TABLE IF NOT EXISTS birthplace (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    city TEXT CONSTRAINT city_length CHECK (LENGTH(city) <= 30),
    region TEXT CONSTRAINT region_length CHECK (LENGTH(region) <= 30),
    country_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT city_region_country_unique UNIQUE (city, region, country_id),
    FOREIGN KEY (country_id) REFERENCES country(id) ON DELETE SET NULL
    );

-- Создание таблицы person
CREATE TABLE IF NOT EXISTS person (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    first_name TEXT CONSTRAINT person_first_name_length CHECK (LENGTH(last_name) <= 30) NOT NULL,
    last_name TEXT CONSTRAINT person_last_name_length CHECK (LENGTH(last_name) <= 30) NOT NULL,
    birth_date TIMESTAMP NOT NULL,
    death_date TIMESTAMP CHECK (death_date > birth_date),
    birthplace_id INT NOT NULL,
    start_career TIMESTAMP,
    end_career TIMESTAMP CHECK (end_career > start_career),
    sex CHAR(1) CONSTRAINT person_gender CHECK (sex='M' OR sex='F') NOT NULL,
    height DECIMAL(10,2) CONSTRAINT height_positive CHECK (height > 100) NOT NULL,
    spouse TEXT CONSTRAINT person_spouse_length CHECK (LENGTH(spouse) <= 150),
    children TEXT CONSTRAINT person_children_length CHECK (LENGTH(children) <= 150),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (birthplace_id) REFERENCES birthplace(id) ON DELETE SET NULL
    );

-- Создание таблицы roles
CREATE TABLE IF NOT EXISTS roles (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    name TEXT CONSTRAINT roles_name_length CHECK (LENGTH(name) <= 50) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Создание таблицы film
CREATE TABLE IF NOT EXISTS film (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    content_id INT NOT NULL UNIQUE,
    year INT CONSTRAINT film_year_positive CHECK (year > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE
    );

-- Создание таблицы content
CREATE TABLE IF NOT EXISTS content (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    title TEXT CONSTRAINT title_length CHECK (LENGTH(title) <= 150) NOT NULL ,
    original_title TEXT CONSTRAINT original_title_length CHECK (LENGTH(original_title) <= 150),
    budget INT CONSTRAINT content_budget_positive CHECK (budget > 0),
    marketing INT CONSTRAINT content_marketing_positive CHECK (marketing > 0),
    premiere TIMESTAMP,
    release TIMESTAMP,
    age_restriction INT CONSTRAINT content_age_restriction_not_negative CHECK (age_restriction >= 0),
    imdb DECIMAL(3,1) CONSTRAINT content_imdb CHECK (imdb > 0 AND imdb <= 10),
    description TEXT CONSTRAINT description_length CHECK (LENGTH(description) <= 1000),
    poster TEXT CONSTRAINT poster_length CHECK (LENGTH(poster) <= 42),
    playback TEXT CONSTRAINT playback_length CHECK (LENGTH(playback) <= 42),
    type CHAR(1) CONSTRAINT type CHECK (type = 'F' OR type = 'S') NOT NULL,
    duration INT CONSTRAINT content_duration_positive CHECK (duration > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Создание таблицы status
CREATE TABLE IF NOT EXISTS status (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    status TEXT CONSTRAINT status_length CHECK (LENGTH(status) <= 15) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Создание таблицы episode
CREATE TABLE IF NOT EXISTS episode (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    season_id INT NOT NULL,
    description TEXT CONSTRAINT episode_description_length CHECK (LENGTH(description) <= 500),
    episode_number INT CONSTRAINT episode_episode_number_positive CHECK (episode_number > 0) NOT NULL,
    viewed BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (season_id) REFERENCES season(id) ON DELETE CASCADE
    );

-- Создание таблицы season
CREATE TABLE IF NOT EXISTS season (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    series_id INT NOT NULL,
    year_start INT CONSTRAINT season_year_start_positive CHECK (year_start > 0),
    year_end INT CHECK (year_end >= year_start) CONSTRAINT season_year_end_positive CHECK (year_end > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE CASCADE
    );

-- Создание таблицы series
CREATE TABLE IF NOT EXISTS series (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    title TEXT CONSTRAINT series_title_length CHECK (LENGTH(title) <= 150) NOT NULL,
    year_start INT CONSTRAINT series_year_start_positive CHECK (year_start > 0),
    year_end INT CHECK (year_end >= year_start) CONSTRAINT series_year_end_positive CHECK (year_end > 0),
    content_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE
    );

-- Создание таблицы users
CREATE TABLE IF NOT EXISTS users (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    name TEXT CONSTRAINT user_name_length CHECK (LENGTH(name) <= 150),
    email TEXT CONSTRAINT email_length CHECK (LENGTH(email) <= 150) UNIQUE NOT NULL,
    password_hashed bytea CONSTRAINT password_hashed_length CHECK (OCTET_LENGTH(password_hashed) <= 32) NOT NULL,
    salt_password bytea CONSTRAINT salt_password_length CHECK (OCTET_LENGTH(salt_password) <= 8) NOT NULL,
    birth_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Создание таблицы review
CREATE TABLE IF NOT EXISTS review (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    user_id INT NOT NULL,
    content_id INT NOT NULL,
    title TEXT CONSTRAINT review_title_length CHECK (LENGTH(title) > 0 AND LENGTH(title) <= 50)NOT NULL,
    text TEXT CONSTRAINT review_text_length CHECK (LENGTH(text) > 0 AND LENGTH(text) <= 1000) NOT NULL,
    content_rating INT CONSTRAINT rating_user_positive CHECK (content_rating > 0 AND content_rating <= 10)  NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL,
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE
    );

-- Создание таблицы nomination
CREATE TABLE IF NOT EXISTS nomination (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    title TEXT CONSTRAINT nomination_title_length CHECK (LENGTH(title) > 0 AND LENGTH(title) <= 50) NOT NULL,
    content_id INT NOT NULL,
    person_id INT,
    award_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE,
    FOREIGN KEY (person_id) REFERENCES person(id) ON DELETE CASCADE ,
    FOREIGN KEY (award_id) REFERENCES award(id) ON DELETE CASCADE
    );

-- Создание таблицы award
CREATE TABLE IF NOT EXISTS award (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    year INT CONSTRAINT award_year_positive CHECK (year > 0),
    name TEXT CONSTRAINT award_name_length CHECK (LENGTH(name) > 0 AND LENGTH(name) <= 50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Создание таблицы rating_of_content
CREATE TABLE IF NOT EXISTS rating_of_content (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    user_id INT NOT NULL,
    value INT CONSTRAINT rating_value_positive CHECK (value > 0 AND value <= 10),
    content_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL,
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE
);

-- Создание таблицы rating_of_person
CREATE TABLE IF NOT EXISTS rating_of_person (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    user_id INT NOT NULL,
    value INT CONSTRAINT rating_value_positive CHECK (value > 0 AND value <= 10),
    person_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    udated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL,
    FOREIGN KEY (person_id) REFERENCES person(id) ON DELETE CASCADE
);

-- Создание таблицы review_likes
CREATE TABLE IF NOT EXISTS review_likes (
    review_id INT NOT NULL,
    user_id INT NOT NULL,
    value BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (review_id) REFERENCES review(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT review_likes_unique UNIQUE (review_id, user_id)
    );

-- Создание таблицы saved_person
CREATE TABLE IF NOT EXISTS saved_person (
    person_id INT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (person_id) REFERENCES person(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT saved_person_unique UNIQUE (person_id, user_id)
    );

-- Создание таблицы content_status
CREATE TABLE IF NOT EXISTS content_status (
    content_id INT NOT NULL,
    user_id INT NOT NULL,
    status_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (status_id) REFERENCES status(id) ON DELETE CASCADE
    );

-- Создание таблицы genre_content
CREATE TABLE IF NOT EXISTS genre_content (
    genre_id INT NOT NULL,
    content_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (genre_id) REFERENCES genre(id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE,
    CONSTRAINT genre_content_unique UNIQUE (genre_id, content_id)
    );

-- Создание таблицы country_content
CREATE TABLE IF NOT EXISTS country_content (
    country_id INT NOT NULL,
    content_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (country_id) REFERENCES country(id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE,
    CONSTRAINT country_content_unique UNIQUE (country_id, content_id)
    );

-- Создание таблицы content_person
CREATE TABLE IF NOT EXISTS content_person (
    content_id INT NOT NULL,
    person_id INT NOT NULL,
    roles_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE,
    FOREIGN KEY (person_id) REFERENCES person(id) ON DELETE CASCADE,
    FOREIGN KEY (roles_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT content_person_unique UNIQUE (content_id, person_id, roles_id)
    );
