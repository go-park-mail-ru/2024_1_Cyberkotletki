```mermaid
erDiagram
    STATIC {
        INT id PK "Уникальный идентификатор"
        STRING name "Имя"
        STRING path "Путь"
        TIMESTAMPTZ created_at "Время создания"
    }
    COUNTRY {
        INT id PK "Уникальный идентификатор страны"
        STRING name "Название страны (AK1)"
    }
    GENRE {
        INT id PK "Уникальный идентификатор жанра"
        STRING name "Название жанра (AK1)"
    }
    PERSON {
        INT id PK "Уникальный идентификатор персоны"
        STRING first_name "Имя персоны"
        STRING last_name "Фамилия персоны"
        DATE birth_date "Дата рождения персоны"
        DATE death_date "Дата смерти персоны"
        DATE start_career "Дата начала карьеры персоны"
        DATE end_career "Дата окончания карьеры персоны"
        CHAR sex "Пол персоны"
        INT height "Рост персоны"
        STRING spouse "Супруг(а) персоны"
        STRING children "Дети персоны"
        INT photo_upload_id FK "Идентификатор загрузки фото"
        TIMESTAMPTZ created_at "Время создания"
        TIMESTAMPTZ updated_at "Время обновления"
    }
    CONTENT {
        INT id PK "Уникальный идентификатор контента"
        STRING title "Название контента"
        STRING original_title "Оригинальное название контента"
        STRING slogan "Слоган контента"
        INT budget "Бюджет контента"
        INT age_restriction "Ограничение по возрасту контента"
        INT audience "Аудитория контента"
        DECIMAL imdb "Рейтинг IMDB контента"
        STRING description "Описание контента"
        INT poster_upload_id FK "Идентификатор загрузки постера"
        INT box_office "Кассовые сборы контента"
        INT marketing_budget "Бюджет маркетинга контента"
        TIMESTAMPTZ created_at "Время создания"
        TIMESTAMPTZ updated_at "Время обновления"
    }
    ROLES {
        INT id PK "Уникальный идентификатор роли"
        STRING name "Название роли (AK1)"
    }
    PERSON_ROLES {
        INT roles_id FK "Идентификатор роли (AK1.1)"
        INT person_id FK "Идентификатор персоны (AK1.2)"
        INT content_id FK "Идентификатор контента (AK1.3)"
    }
    MOVIE {
        INT id PK "Уникальный идентификатор фильма"
        INT content_id FK "Идентификатор контента (AK1)"
        DATE premiere "Дата премьеры фильма"
        DATE release "Дата выпуска фильма"
        INT duration "Продолжительность фильма"
    }
    SERIES {
        INT id PK "Уникальный идентификатор сериала"
        INT content_id FK "Идентификатор контента (AK1)"
        INT year_start "Год начала сериала"
        INT year_end "Год окончания сериала"
    }
    SEASON {
        INT id PK "Уникальный идентификатор сезона"
        INT series_id FK "Идентификатор сериала"
        INT year_start "Год начала сезона"
        INT year_end "Год окончания сезона"
    }
    EPISODE {
        INT id PK "Уникальный идентификатор эпизода"
        INT season_id FK "Идентификатор сезона"
        INT episode_number "Номер эпизода"
        STRING title "Название эпизода"
    }
    USERS {
        INT id PK "Уникальный идентификатор пользователя"
        STRING name "Имя пользователя"
        STRING email "Email пользователя (AK1)"
        BYTEA password_hashed "Хешированный пароль пользователя"
        BYTEA salt_password "Соль для пароля пользователя"
        INT avatar_upload_id FK "Идентификатор загрузки аватара"
        TIMESTAMPTZ created_at "Время создания"
        TIMESTAMPTZ updated_at "Время обновления"
    }
    REVIEW {
        INT id PK "Уникальный идентификатор обзора"
        INT user_id FK "Идентификатор пользователя (AK1.1)"
        INT content_id FK "Идентификатор контента (AK1.2)"
        STRING title "Название обзора"
        STRING text "Текст обзора"
        INT content_rating "Рейтинг контента"
        TIMESTAMPTZ created_at "Время создания"
        TIMESTAMPTZ updated_at "Время обновления"
    }
    COMPILATION_TYPE {
        INT id PK "Уникальный идентификатор типа подборки"
        STRING type "Тип подборки (AK1)"
    }
    COMPILATION {
        INT id PK "Уникальный идентификатор подборки"
        STRING title "Название подборки"
        INT compilation_type_id FK "Идентификатор типа подборки"
        INT poster_upload_id FK "Идентификатор загрузки постера"
    }
    COMPILATION_CONTENT {
        INT compilation_id FK "Идентификатор подборки (AK1.1)"
        INT content_id FK "Идентификатор контента (AK1.2)"
    }
    REVIEW_LIKE {
        INT review_id FK "Идентификатор обзора  (AK1.1)"
        INT user_id FK "Идентификатор пользователя (AK1.2)"
        BOOLEAN value "Значение лайка"
        TIMESTAMPTZ updated_at "Время обновления"
    }
    GENRE_CONTENT {
        INT genre_id FK "Идентификатор жанра (AK1.1)"
        INT content_id FK "Идентификатор контента (AK1.2)"
    }
    COUNTRY_CONTENT {
        INT country_id FK "Идентификатор страны (AK1.1)"
        INT content_id FK "Идентификатор контента (AK1.2)"
    }
    CONTENT_TYPE {
        INT id PK "Уникальный идентификатор типа контента"
        STRING type "Тип контента"
    }
    CONTENT_CONTENT_TYPE {
        INT content_id FK "Идентификатор контента (AK1.1)"
        INT content_type_id FK "Идентификатор типа контента (AK1.2)"
    }
    
    STATIC ||--o{ PERSON : "photo_upload_id"
    STATIC ||--o{ CONTENT : "poster_upload_id"
    STATIC ||--o{ USERS : "avatar_upload_id"
    STATIC ||--o{ COMPILATION : "poster_upload_id"
    COUNTRY ||--o{ COUNTRY_CONTENT : "country_id"
    GENRE ||--o{ GENRE_CONTENT : "genre_id"
    PERSON ||--o{ PERSON_ROLES : "person_id"
    PERSON ||--o{ USERS : "id"
    PERSON ||--o{ REVIEW : "user_id"
    PERSON ||--o{ SAVED_PERSON : "person_id"
    CONTENT ||--o{ PERSON_ROLES : "content_id"
    CONTENT ||--o{ CONTENT_STATUS : "content_id"
    CONTENT ||--o{ REVIEW : "content_id"
    CONTENT ||--o{ COMPILATION_CONTENT : "content_id"
    CONTENT ||--o{ GENRE_CONTENT : "content_id"
    CONTENT ||--o{ COUNTRY_CONTENT : "content_id"
    CONTENT ||--o{ CONTENT_CONTENT_TYPE : "content_id"
    CONTENT_TYPE ||--o{ CONTENT_CONTENT_TYPE : "content_type_id"
    CONTENT ||--o{ MOVIE : "content_id"
    CONTENT ||--o{ SERIES : "content_id"
    ROLES ||--o{ PERSON_ROLES : "roles_id"
    USERS ||--o{ CONTENT_STATUS : "users_id"
    USERS ||--o{ REVIEW : "user_id"
    USERS ||--o{ SAVED_PERSON : "users_id"
    USERS ||--o{ REVIEW_LIKE : "user_id"
    COMPILATION_TYPE ||--o{ COMPILATION : "compilation_type_id"
    REVIEW ||--o{ REVIEW_LIKE : "review_id"
    SERIES ||--o{ SEASON : "series_id"
    SEASON ||--o{ EPISODE : "season_id"