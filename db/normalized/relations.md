## Таблица STATIC

Таблица `STATIC` содержит информацию о статических данных, таких как изображения, которые загружаются на сайт. 
Это могут быть аватары пользователей, постеры контента и т.д.

<p> Функциональные зависимости: </p>

- `{id} -> {name, path, created_at}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты id, name, path, created_at являются атомарными.
- 2 НФ: Атрибуты name, path, created_at полностью функционально зависят от первичного ключа id.
- 3 НФ: Атрибуты name, path, created_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    STATIC {
        INT id PK "Уникальный идентификатор"
        STRING name "Имя"
        STRING path "Путь"
        TIMESTAMPTZ created_at "Время создания"
    }
```

## Таблица COUNTRY

Таблица `COUNTRY` содержит информацию о странах.

<p> Функциональные зависимости: </p>

- `{id} -> {name}`
- `{name} -> {id}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты id, name являются атомарными.
- 2 НФ: Атрибут name полностью функционально зависит от первичного ключа id.
- 3 НФ: Атрибут name не зависит от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    COUNTRY {
        INT id PK "Уникальный идентификатор страны"
        STRING name "Название страны (AK1)"
    }
```

## Таблица GENRE

Таблица `GENRE` содержит информацию о жанрах контента.

<p> Функциональные зависимости: </p>

- `{id} -> {name}`
- `{name} -> {id}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты id, name являются атомарными.
- 2 НФ: Атрибут name полностью функционально зависит от первичного ключа id.
- 3 НФ: Атрибут name не зависит от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    GENRE {
        INT id PK "Уникальный идентификатор жанра"
        STRING name "Название жанра (AK1)"
    }
```

## Таблица PERSON

Таблица `PERSON` содержит информацию о персонах, участвующих в создании контента.

<p> Функциональные зависимости: </p>

- `{id} -> {first_name, last_name, birth_date, death_date, start_career, end_career, sex, height, spouse, children, photo_upload_id, created_at, updated_at}`

<p> Нормальные формы: <p>

- 1 НФ: Все атрибуты являются атомарными.
- 2 НФ: Все атрибуты полностью функционально зависят от первичного ключа id.
- 3 НФ: Все атрибуты не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
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
```

## Таблица CONTENT

Таблица `CONTENT` содержит информацию о контенте.

<p> Функциональные зависимости: </p>

- `{id} -> {title, original_title, slogan, budget, age_restriction, audience, imdb, description, poster_upload_id, box_office, marketing_budget, created_at, updated_at}`

<p> Нормальные формы: <p>

- 1 НФ: Все атрибуты являются атомарными.
- 2 НФ: Все атрибуты полностью функционально зависят от первичного ключа id.
- 3 НФ: Все атрибуты не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
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
```

## Таблица ROLE_DATA

Таблица `ROLE_DATA` содержит информацию о ролях персон в контенте. Например, актер, режиссер, продюсер и т.д.

<p> Функциональные зависимости: </p>

- `{id} -> {name}`
- `{name} -> {id}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты id, name являются атомарными.
- 2 НФ: Атрибут name полностью функционально зависит от первичного ключа id.
- 3 НФ: Атрибут name не зависит от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    ROLE_DATA {
        INT id PK "Уникальный идентификатор роли"
        STRING name "Название роли (AK1)"
    }
```

## Таблица PERSON_ROLE

Таблица `PERSON_ROLE` содержит информацию о связи между персонами, ролями и контентом.

<p> Функциональные зависимости:  - </p>

<p> Нормальные формы: <p>

- 1 НФ: Все атрибуты являются атомарными.
- 2 НФ: Нет атрибутов, которые зависят от части составного ключа.
- 3 НФ: Нет атрибутов, которые зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют транзитивные зависимости.

```mermaid
erDiagram
    PERSON_ROLE {
        INT role_id FK "Идентификатор роли (AK1.1)"
        INT person_id FK "Идентификатор персоны (AK1.2)"
        INT content_id FK "Идентификатор контента (AK1.3)"
    }
```

## Таблица MOVIE

Таблица `MOVIE` содержит информацию о фильмах.

<p> Функциональные зависимости: </p>

- `{id} -> {content_id, premiere, release, duration}`
- `{content_id} -> {id, premiere, release, duration}`

<p> Нормальные формы: <p>

- 1 НФ: Все атрибуты являются атомарными.
- 2 НФ: Все атрибуты полностью функционально зависят от первичного ключа id.
- 3 НФ: Все атрибуты не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    MOVIE {
        INT id PK "Уникальный идентификатор фильма"
        INT content_id FK "Идентификатор контента (AK1)"
        DATE premiere "Дата премьеры фильма"
        DATE release "Дата выпуска фильма"
        INT duration "Продолжительность фильма"
    }
```

## Таблица SERIES

Таблица `TV_SERIES` содержит информацию о сериалах.

<p> Функциональные зависимости: </p>

- `{id} -> {year_start, year_end, content_id}`
- `{content_id} -> {id, year_start, year_end}`

<p> Нормальные формы: <p>

- 1 НФ: Все атрибуты являются атомарными.
- 2 НФ: Все атрибуты полностью функционально зависят от первичного ключа id.
- 3 НФ: Все атрибуты не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    SERIES {
        INT id PK "Уникальный идентификатор сериала"
        INT year_start "Год начала сериала"
        INT year_end "Год окончания сериала"
        INT content_id FK "Идентификатор контента"
    }
```

## Таблица SEASON

Таблица `SEASON` содержит информацию о сезонах сериалов.

<p> Функциональные зависимости: </p>

- `{id} -> {tv_series_id, title, year_start, year_end}`

<p> Нормальные формы: <p>

- 1 НФ: Все атрибуты являются атомарными.
- 2 НФ: Все атрибуты полностью функционально зависят от первичного ключа id.
- 3 НФ: Все атрибуты не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    SEASON {
        INT id PK "Уникальный идентификатор сезона"
        INT tv_series_id FK "Идентификатор сериала (АК1.2)"
        STRING title "Название сезона (АК1.2)"
        INT year_start "Год начала сезона"
        INT year_end "Год окончания сезона"
    }
```

## Таблица EPISODE

Таблица `EPISODE` содержит информацию о сериях сериалов.

<p> Функциональные зависимости: </p>

- `{id} -> {season_id, episode_number, title}`

<p> Нормальные формы: <p>

- 1 НФ: Все атрибуты являются атомарными.
- 2 НФ: Все атрибуты полностью функционально зависят от первичного ключа id.
- 3 НФ: Все атрибуты не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    EPISODE {
        INT id PK "Уникальный идентификатор эпизода"
        INT season_id FK "Идентификатор сезона"
        INT episode_number "Номер эпизода"
        STRING title "Название эпизода"
    }
```


## Таблица USER_DATA

Таблица `USER_DATA` содержит информацию о пользователях.

<p> Функциональные зависимости: </p>

- `{id} -> {name, email, password_hashed, salt_password, avatar_upload_id, created_at, updated_at}`
- `{email} -> {id, name, password_hashed, salt_password, avatar_upload_id, created_at, updated_at}`

<p> Нормальные формы: <p>

- 1 НФ: Все атрибуты являются атомарными.
- 2 НФ: Все атрибуты полностью функционально зависят от первичного ключа id.
- 3 НФ: Все атрибуты не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    USER_DATA {
        INT id PK "Уникальный идентификатор пользователя"
        STRING name "Имя пользователя"
        STRING email "Email пользователя (AK1)"
        BYTEA password_hashed "Хешированный пароль пользователя"
        BYTEA salt_password "Соль для пароля пользователя"
        INT avatar_upload_id FK "Идентификатор загрузки аватара"
        TIMESTAMPTZ created_at "Время создания"
        TIMESTAMPTZ updated_at "Время обновления"
    }
```

## Таблица REVIEW

Таблица `REVIEW` содержит информацию о ревью на контент, которые делают пользователи.

<p> Функциональные зависимости: </p>

- `{id} -> {user_id, content_id, title, text, content_rating, created_at, updated_at}`
- `{user_id, content_id} -> {id, title, text, content_rating, created_at, updated_at}`

<p> Нормальные формы: <p>

- 1 НФ: Все атрибуты являются атомарными.
- 2 НФ: Все атрибуты полностью функционально зависят от первичного ключа id.
- 3 НФ: Все атрибуты не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
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
```

## Таблица COMPILATION_TYPE

Таблица `COMPILATION_TYPE` содержит информацию о типах подборок контента. Например, годы, фильмы, сериалы и т.д.

<p> Функциональные зависимости: </p>

- `{id} -> {type}`
- `{type} -> {id}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты id, type являются атомарными.
- 2 НФ: Атрибут type полностью функционально зависит от первичного ключа id.
- 3 НФ: Атрибут type не зависит от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    COMPILATION_TYPE {
        INT id PK "Уникальный идентификатор типа подборки"
        STRING type "Тип подборки (AK1)"
    }
```

## Таблица COMPILATION

Таблица `COMPILATION` содержит информацию о подборках контента. Например, подборка лучших фильмов 2021 года.

<p> Функциональные зависимости: </p>

- `{id} -> {title, compilation_type_id, poster_upload_id}`

<p> Нормальные формы: <p>

- 1 НФ: Все атрибуты являются атомарными.
- 2 НФ: Все атрибуты полностью функционально зависят от первичного ключа id.
- 3 НФ: Все атрибуты не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    COMPILATION {
        INT id PK "Уникальный идентификатор подборки"
        STRING title "Название подборки"
        INT compilation_type_id FK "Идентификатор типа подборки"
        INT poster_upload_id FK "Идентификатор загрузки постера"
    }
```

## Таблица COMPILATION_CONTENT

Таблица `COMPILATION_CONTENT` содержит информацию о связи между подборками и контентом.

<p> Функциональные зависимости: - </p>

<p> Нормальные формы: <p>

- 1 НФ: Все атрибуты являются атомарными.
- 2 НФ: Нет атрибутов, которые зависят от части составного ключа.
- 3 НФ: Нет атрибутов, которые зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют транзитивные зависимости.

```mermaid
erDiagram
    COMPILATION_CONTENT {
        INT compilation_id FK "Идентификатор подборки (PK1.1)"
        INT content_id FK "Идентификатор контента (PK1.2)"
    }
```

## Таблица REVIEW_LIKE

Таблица `REVIEW_LIKE` содержит информацию о лайках, которые пользователи ставят ревью.

<p> Функциональные зависимости: </p>

- `{review_id, user_id} -> {value, updated_at}`

<p> Нормальные формы: <p>

- 1 НФ: Все атрибуты являются атомарными.
- 2 НФ: Атрибуты value, updated_at полностью функционально зависят от составного ключа {review_id, user_id}.
- 3 НФ: Атрибуты value, updated_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют транзитивные зависимости.

```mermaid
erDiagram
    REVIEW_LIKE {
        INT review_id FK "Идентификатор обзора  (PK1.1)"
        INT user_id FK "Идентификатор пользователя (PK1.2)"
        BOOLEAN value "Значение лайка"
        TIMESTAMPTZ updated_at "Время обновления"
    }
```

## Таблица GENRE_CONTENT

Таблица `GENRE_CONTENT` содержит информацию о связи между жанрами и контентом.

<p> Функциональные зависимости: - </p>

<p> Нормальные формы: <p>

- 1 НФ: Все атрибуты являются атомарными.
- 2 НФ: Нет атрибутов, которые зависят от части составного ключа.
- 3 НФ: Нет атрибутов, которые зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют транзитивные зависимости.

```mermaid
erDiagram
    GENRE_CONTENT {
        INT genre_id FK "Идентификатор жанра (PK1.1)"
        INT content_id FK "Идентификатор контента (PK1.2)"
    }
```

## Таблица COUNTRY_CONTENT

Таблица `COUNTRY_CONTENT` содержит информацию о связи между странами и контентом.

<p> Функциональные зависимости: - </p>

<p> Нормальные формы: <p>

- 1 НФ: Все атрибуты являются атомарными.
- 2 НФ: Нет атрибутов, которые зависят от части составного ключа.
- 3 НФ: Нет атрибутов, которые зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют транзитивные зависимости.

```mermaid
erDiagram
    COUNTRY_CONTENT {
        INT country_id FK "Идентификатор страны (PK1.1)"
        INT content_id FK "Идентификатор контента (PK1.2)"
    }
```

## Таблица CONTENT_TYPE

Таблица `CONTENT_TYPE` содержит информацию о типах контента: фильм, сериал.
<p> Функциональные зависимости: </p>

- `{id} -> {type}`

<p> Нормальные формы: <p>

- 1 НФ: Все атрибуты являются атомарными.
- 2 НФ: Все атрибуты полностью функционально зависят от первичного ключа id.
- 3 НФ: Все атрибуты не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    CONTENT_TYPE {
        INT id PK "Уникальный идентификатор типа контента"
        STRING type "Тип контента AK1"
    }
```

## Таблица CONTENT_CONTENT_TYPE

Таблица `CONTENT_CONTENT_TYPE` содержит информацию о связи между контентом и его типом.

<p> Функциональные зависимости: - </p>

<p> Нормальные формы: <p>

- 1 НФ: Все атрибуты являются атомарными.
- 2 НФ: Нет атрибутов, которые зависят от части составного ключа.
- 3 НФ: Нет атрибутов, которые зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют транзитивные зависимости.

```mermaid
erDiagram
    CONTENT_CONTENT_TYPE {
        INT content_id FK "Идентификатор контента (PK1.1)"
        INT content_type_id FK "Идентификатор типа контента (PK1.2)"
    }
```