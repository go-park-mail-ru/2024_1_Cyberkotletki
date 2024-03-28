# Основные таблицы

## Таблица AUDIENCE

Таблица `AUDIENCE` содержит информацию об аудитории контента в определенной стране. 
У одного контента может быть несколько аудиторий в разных странах.
<p> Функциональные зависимости: <p> 

- `{ID} -> {CONTENT_ID, SIZE_IN_THOUSANDS, COUNTRY_ID, created_at}`
- `{CONTENT_ID} -> CONTENT {ID}`

<p> Нормальные формы: <p> 

- 1 НФ: Атрибуты ID, CONTENT_ID, SIZE_IN_THOUSANDS, COUNTRY_ID, created_at являются атомарными.
- 2 НФ: Атрибуты CONTENT_ID, SIZE_IN_THOUSANDS, COUNTRY_ID, created_at полностью функционально зависят от первичного ключа ID. 
- 3 НФ: Атрибуты CONTENT_ID, SIZE_IN_THOUSANDS, COUNTRY_ID, created_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.


```mermaid
erDiagram
    AUDIENCE {
        INT ID PK "Уникальный идентификатор аудитории"
        INT CONTENT_ID FK "Идентификатор контента, который имеет эту аудиторию"
        INT SIZE_IN_THOUSANDS "Размер аудитории в тысячах"
        INT COUNTRY_ID "Идентификатор страны аудитории"
        TIME created_at "время создания кортежа, для отладки"
    }
```
## Таблица COUNTRY

Таблица `COUNTRY` содержит информацию о странах, в которых был снят/произведен контент.
<p> Функциональные зависимости: </p>

- `{ID} -> {NAME, created_at}`

<p> Нормальные формы: <p> 

- 1 НФ: Атрибуты ID, NAME, created_at являются атомарными.
- 2 НФ: Атрибуты NAME, created_at полностью функционально зависят от первичного ключа ID.
- 3 НФ: Атрибуты NAME, created_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.


```mermaid
erDiagram
    COUNTRY {
        INT ID PK "Уникальный идентификатор страны"
        STRING NAME "Название страны"
        TIME created_at "время создания кортежа, для отладки"
    }
```
## Таблица BOXOFFICE

Таблица `BOXOFFICE` содержит информацию о кассовых сборах определенной суммы в определенной стране.
У одного контента может быть несколько кассовых сборов в разных странах.
<p> Функциональные зависимости: </p>

- `{ID} -> {CONTENT_ID, COUNTRY_ID, REVENUE, created_at}`
- `{CONTENT_ID} -> CONTENT {ID}`
- `{COUNTRY_ID} -> COUNTRY {ID}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты ID, CONTENT_ID, COUNTRY_ID, REVENUE, created_at являются атомарными.
- 2 НФ: Атрибуты CONTENT_ID, COUNTRY_ID, REVENUE, created_at полностью функционально зависят от первичного ключа ID.
- 3 НФ: Атрибуты CONTENT_ID, COUNTRY_ID, REVENUE, created_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    BOXOFFICE {
        INT ID PK "Уникальный идентификатор кассовых сборов"
        INT CONTENT_ID FK "Идентификатор контента, который имеет этот кассовый сбор"
        INT COUNTRY_ID "Идентификатор страны кассовых сборов"
        INT REVENUE "Сумма сборов"
        TIME created_at "время создания кортежа, для отладки"
    }
```

## Таблица GENRE

Таблица `GENRE` содержит информацию о жанрах контента на русском и английском языках.
<p> Функциональные зависимости: </p>

- `{ID} -> {NAME, NAME_RU, created_at}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты ID, NAME, NAME_RU, created_at являются атомарными.
- 2 НФ: Атрибуты NAME, NAME_RU, created_at полностью функционально зависят от первичного ключа ID.
- 3 НФ: Атрибуты NAME, NAME_RU, created_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    GENRE {
        INT ID PK "Уникальный идентификатор жанра"
        STRING NAME "Название жанра на английском"
        STRING NAME_RU "Название жанра на русском"
        TIME created_at "время создания кортежа, для отладки"
    }
```

## Таблица BIRTHPLACE

Таблица `BIRTHPLACE` содержит информацию о месте рождения персоны.
<p> Функциональные зависимости: </p>

- `{ID} -> {CITY, REGION, COUNTRY_ID, created_at}`
- `{COUNTRY_ID} -> COUNTRY {ID}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты ID, CITY, REGION, COUNTRY_ID, created_at являются атомарными.
- 2 НФ: Атрибуты CITY, REGION, COUNTRY_ID, created_at полностью функционально зависят от первичного ключа ID.
- 3 НФ: Атрибуты CITY, REGION, COUNTRY_ID, created_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    BIRTHPLACE{
        INT ID PK "Уникальный идентификатор места рождения"
        STRING CITY "Город рождения"
        STRING REGION "Регион рождения"
        INT COUNTRY_ID "Идентификатор страны рождения"
        TIME created_at "время создания кортежа, для отладки"
    }
```

## Таблица PERSON

Таблица `PERSON` содержит информацию о персонах, которые участвовали в создании контента.
<p> Функциональные зависимости: </p>

- `{ID} -> {FIRST_NAME, LAST_NAME, BIRTH_DATE, BIRTHPLACE_ID, DEATH_DATE, START_CAREER, END_CAREER, SEX, 
    PHOTO, HEIGHT, SPOUSE, CHILDREN, created_at, updated_at}`
- `{BIRTHPLACE_ID} -> BIRTHPLACE {ID}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты ID, FIRST_NAME, LAST_NAME, BIRTH_DATE, BIRTHPLACE_ID, DEATH_DATE, START_CAREER, END_CAREER, SEX, 
    PHOTO, HEIGHT, SPOUSE, CHILDREN, created_at, updated_at являются атомарными.
- 2 НФ: Атрибуты FIRST_NAME, LAST_NAME, BIRTH_DATE, BIRTHPLACE_ID, DEATH_DATE, START_CAREER, END_CAREER, SEX, 
    PHOTO, HEIGHT, SPOUSE, CHILDREN, created_at, updated_at полностью функционально зависят от первичного ключа ID.
- 3 НФ: Атрибуты FIRST_NAME, LAST_NAME, BIRTH_DATE, BIRTHPLACE_ID, DEATH_DATE, START_CAREER, END_CAREER, SEX, 
    PHOTO, HEIGHT, SPOUSE, CHILDREN, created_at, updated_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    PERSON {
        INT ID PK "Уникальный идентификатор персоны"
        STRING FIRST_NAME "Имя персоны"
        STRING LAST_NAME "Фамилия персоны"
        TIME BIRTH_DATE "Дата рождения персоны"
        INT BIRTHPLACE_ID FK "Идентификатор места рождения персоны"
        TIME DEATH_DATE "Дата смерти персоны"
        TIME START_CAREER "Дата начала карьеры персоны"
        TIME END_CAREER "Дата окончания карьеры персоны"
        STRING SEX "Пол персоны"
        STRING PHOTO "Фотография персоны"
        INT HEIGHT "Рост персоны"
        STRING SPOUSE "Супруг(а) персоны"
        STRING CHILDREN "Дети персоны"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
```

## Таблица ROLES

Таблица `ROLES` содержит информацию о ролях персоны в контенте, наприер, актер, режиссер, дублер и т.д.
<p> Функциональные зависимости: <p> 

- `{ID} -> {NAME, created_at}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты ID, NAME, created_at являются атомарными.
- 2 НФ: Атрибуты NAME, created_at полностью функционально зависят от первичного ключа ID.
- 3 НФ: Атрибуты NAME, created_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    ROLES{
        INT ID PK "Уникальный идентификатор роли персоны в контенте"
        STRING NAME "Название роли (актер, режиссер, дублер и т.д.)"
        TIME created_at "время создания кортежа, для отладки"
    }
```

## Таблица FILM

Таблица `FILM` содержит информацию о фильмах. В фильме есть данные контента и некоторые специфичные для фильма данные.
<p> Функциональные зависимости: <p> 

- `{ID} -> {CONTENT_ID, YEAR, created_at, updated_at}`
- `{CONTENT_ID} -> CONTENT {ID}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты ID, CONTENT_ID, YEAR, created_at, updated_at являются атомарными.
- 2 НФ: Атрибуты CONTENT_ID, YEAR, created_at, updated_at полностью функционально зависят от первичного ключа ID.
- 3 НФ: Атрибуты CONTENT_ID, YEAR, created_at, updated_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    FILM {
        INT ID PK "Уникальный идентификатор фильма"
        INT CONTENT_ID FK "Идентификатор контента, который относится к фильму"
        INT YEAR "Год выпуска фильма"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
```

## Таблица CONTENT

Таблица `CONTENT` содержит информацию о контенте, который может быть фильмом или сериалом. Контент собирает в себе общую 
информацию, которая может быть общей как для фильмов, так и для сериалов.
<p> Функциональные зависимости: <p> 

- `{ID} -> {TITLE, ORIGINAL_TITLE, BUDGET, MARKETING_BUDGET, PREMIERE, RELEASE, AGE_RESTRICTION, IMDB, DESCRIPTION, POSTER, 
    PLAYBACK, DURATION, created_at, updated_at}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты ID, TITLE, ORIGINAL_TITLE, BUDGET, MARKETING_BUDGET, PREMIERE, RELEASE, AGE_RESTRICTION, IMDB, DESCRIPTION, POSTER, 
    PLAYBACK, DURATION, created_at, updated_at являются атомарными.
- 2 НФ: Атрибуты TITLE, ORIGINAL_TITLE, BUDGET, MARKETING_BUDGET, PREMIERE, RELEASE, AGE_RESTRICTION, IMDB, DESCRIPTION, POSTER, 
    PLAYBACK, DURATION, created_at, updated_at полностью функционально зависят от первичного ключа ID.
- 3 НФ: Атрибуты TITLE, ORIGINAL_TITLE, BUDGET, MARKETING_BUDGET, PREMIERE, RELEASE, AGE_RESTRICTION, IMDB, DESCRIPTION, POSTER, 
    PLAYBACK, DURATION, created_at, updated_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    CONTENT {
        INT ID PK "Уникальный идентификатор контента"
        STRING TITLE "Название контента"
        STRING ORIGINAL_TITLE "Оригинальное название контента"
        INT BUDGET "Бюджет контента"
        INT MARKETING_BUDGET "Бюджет маркетинга контента"
        TIME PREMIERE "Дата премьеры контента"
        TIME RELEASE "Дата выпуска контента"
        INT AGE_RESTRICTION "Возрастное ограничение контента"
        INT IMDB "Рейтинг IMDB контента"
        STRING DESCRIPTION "Описание контента"
        STRING POSTER "Постер контента"
        STRING PLAYBACK "Воспроизведение на заднем плане небольшого фрагмента видео контента"
        INT DURATION "Продолжительность контента в минутах"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
```

## Таблица STATUS

Таблица `STATUS` содержит информацию о статусе контента у пользователя. Статус может быть просмотрен, запланирован, пересматривается, добавлен в избранное.
<p> Функциональные зависимости: <p> 

- `{ID} -> {STATUS, created_at,}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты ID, STATUS, created_at являются атомарными.
- 2 НФ: Атрибуты STATUS, created_at полностью функционально зависят от первичного ключа ID.
- 3 НФ: Атрибуты STATUS, created_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    STATUS{
        INT ID PK "Уникальный идентификатор статуса"
        STRING STATUS "Статус контента (Viewed, Planned, Reconsidering, Favourites(избранное))"
        TIME created_at "временная метка создания кортежа"
    }
```

## Таблица EPISODE

Таблица `EPISODE` содержит информацию об эпизодах сериала. В сезоне может быть несколько эпизодов.
<p> Функциональные зависимости: <p> 

- `{ID} -> {SEASON_ID, DESCRIPTION, EPISODE_NUMBER, VIEWED, created_at, updated_at}`
- `{SEASON_ID} -> SEASON {ID}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты ID, SEASON_ID, DESCRIPTION, EPISODE_NUMBER, VIEWED, created_at, updated_at являются атомарными.
- 2 НФ: Атрибуты SEASON_ID, DESCRIPTION, EPISODE_NUMBER, VIEWED, created_at, updated_at полностью функционально зависят от первичного ключа ID.
- 3 НФ: Атрибуты SEASON_ID, DESCRIPTION, EPISODE_NUMBER, VIEWED, created_at, updated_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    EPISODE {
        INT ID PK "Уникальный идентификатор эпизода"
        INT SEASON_ID FK "Идентификатор сезона, к которому относится эпизод"
        STRING DESCRIPTION "Описание эпизода"
        INT EPISODE_NUMBER "Номер эпизода"
        STRING VIEWED "Просмотрен ли эпизод"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
```

## Таблица SEASON

Таблица `SEASON` содержит информацию о сезонах сериала. В сериале может быть несколько сезонов.
<p> Функциональные зависимости: <p> 

- `{ID} -> {SERIES_ID, YEAR_START, YEAR_END, created_at, updated_at}`
- `{SERIES_ID} -> SERIES {ID}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты ID, SERIES_ID, YEAR_START, YEAR_END, created_at, updated_at являются атомарными.
- 2 НФ: Атрибуты SERIES_ID, YEAR_START, YEAR_END, created_at, updated_at полностью функционально зависят от первичного ключа ID.
- 3 НФ: Атрибуты SERIES_ID, YEAR_START, YEAR_END, created_at, updated_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    SEASON {
        INT ID PK "Уникальный идентификатор сезона"
        INT SERIES_ID FK "Идентификатор сериала, к которому относится сезон"
        INT YEAR_START "Год начала сезона"
        INT YEAR_END "Год окончания сезона"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
```

## Таблица SERIES

Таблица `SERIES` содержит информацию о сериалах. Сериал содержит данные контента и некоторые специфичные для сериала данные.
<p> Функциональные зависимости: <p> 

- `{ID} -> {CONTENT_ID, TITLE, YEAR_START, YEAR_END, created_at, updated_at}`
- `{CONTENT_ID} -> CONTENT {ID}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты ID, CONTENT_ID, TITLE, YEAR_START, YEAR_END, created_at, updated_at являются атомарными.
- 2 НФ: Атрибуты CONTENT_ID, TITLE, YEAR_START, YEAR_END, created_at, updated_at полностью функционально зависят от первичного ключа ID.
- 3 НФ: Атрибуты CONTENT_ID, TITLE, YEAR_START, YEAR_END, created_at, updated_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    SERIES {
        INT ID PK "Уникальный идентификатор сериала"
        INT CONTENT_ID FK "Идентификатор контента, который относится к  сериалу"
        STRING TITLE "Название сериала"
        INT YEAR_START "Год начала сериала"
        INT YEAR_END "Год окончания сериала"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
```

## Таблица USERS

Таблица `USERS` содержит информацию о пользователях. Пользователь может оставлять комментарии, ставить оценки, сохранять контент и персон.
<p> Функциональные зависимости: <p> 

- `{ID} -> {NAME, EMAIL, PASSWORD_HASHED, SALT_PASSWORD, BIRTH_DATE, created_at, updated_at}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты ID, NAME, EMAIL, PASSWORD_HASHED, SALT_PASSWORD, BIRTH_DATE, created_at, updated_at являются атомарными.
- 2 НФ: Атрибуты NAME, EMAIL, PASSWORD_HASHED, SALT_PASSWORD, BIRTH_DATE, created_at, updated_at полностью функционально зависят от первичного ключа ID.
- 3 НФ: Атрибуты NAME, EMAIL, PASSWORD_HASHED, SALT_PASSWORD, BIRTH_DATE, created_at, updated_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    USERS {
        INT ID PK "Уникальный идентификатор пользователя"
        STRING NAME "Имя пользователя"
        STRING EMAIL "Электронная почта пользователя"
        STRING PASSWORD_HASHED "Хэш пароля пользователя"
        STRING SALT_PASSWORD "Соль для генерации хэша пароля"
        TIME BIRTH_DATE "День рождения пользователя"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
```

## Таблица REVIEW

Таблица `REVIEW` содержит информацию о рецензиях на контент. Пользователь может оставить рецензию на контент.
<p> Функциональные зависимости: <p> 

- `{ID} -> {USERS_ID, CONTENT_ID, TITLE, TEXT, RATING_UUSERS, created_at, updated_at}`
- `{USERS_ID} -> USERS {ID}`
- `{CONTENT_ID} -> CONTENT {ID}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты ID, USERS_ID, CONTENT_ID, TITLE, TEXT, RATING_USERS, created_at, updated_at являются атомарными.
- 2 НФ: Атрибуты USERS_ID, CONTENT_ID, TITLE, TEXT, RATING_USERS, created_at, updated_at полностью функционально зависят от первичного ключа ID.
- 3 НФ: Атрибуты USERS_ID, CONTENT_ID, TITLE, TEXT, RATING_USERS, created_at, updated_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    REVIEW {
        INT ID PK "Уникальный идентификатор комментария"
        INT USERS_ID FK "Идентификатор пользователя, оставившего комментарий"
        INT CONTENT_ID FK "Идентификатор контента, к которому относится комментарий"
        STRING TITLE "Заголовок комментария"
        STRING TEXT "Текст комментария"
        INT RATING_UUSERS "Оценка контента пользователем, оставившим комментарий"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
```

## Таблица NOMINATION

Таблица `NOMINATION` содержит информацию о номинациях. У награды может быть несколько номинаций. 
Номинацию получают либо персона (за контент), либо контент.
<p> Функциональные зависимости: <p> 

- `{ID} -> {TITLE, CONTENT_ID, PERSON_ID, AWARD_ID, created_at}`
- `{CONTENT_ID} -> CONTENT {ID}`
- `{PERSON_ID} -> PERSON {ID}`
- `{AWARD_ID} -> AWARD {ID}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты ID, TITLE, CONTENT_ID, PERSON_ID, AWARD_ID, created_at являются атомарными.
- 2 НФ: Атрибуты TITLE, CONTENT_ID, PERSON_ID, AWARD_ID, created_at полностью функционально зависят от первичного ключа ID.
- 3 НФ: Атрибуты TITLE, CONTENT_ID, PERSON_ID, AWARD_ID, created_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    NOMINATION {
        INT ID PK "Уникальный идентификатор номинации"
        STRING TITLE "Название номинации"
        INT CONTENT_ID FK "Идентификатор контента, за который дана номинация"
        INT PERSON_ID FK  "Идентификатор персоны, которой присвоена номинация"
        INT AWARD_ID FK  "Идентификатор награды"
        TIME created_at "время создания кортежа, для отладки"
    }
```

## Таблица AWARD

Таблица `AWARD` содержит информацию о наградах. 
<p> Функциональные зависимости: <p> 

- `{ID} -> {YEAR, NAME, created_at}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты ID, YEAR, NAME, created_at являются атомарными.
- 2 НФ: Атрибуты YEAR, NAME, created_at полностью функционально зависят от первичного ключа ID.
- 3 НФ: Атрибуты YEAR, NAME, created_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    AWARD {
        INT ID PK "Уникальный идентификатор награды"
        INT YEAR "Год присуждения награды"
        STRING NAME "Тип награды"
        TIME created_at "время создания кортежа, для отладки"
    }
```

## Таблица RATING_OF_CONTENT

Таблица `RATING_OF_CONTENT` содержит информацию о рейтингах контента.
<p> Функциональные зависимости: <p> 

- `{ID} -> {CONTENT_ID, USERS_ID, VALUE, created_at}`
- `{CONTENT_ID} -> CONTENT {ID}`
- `{USERS_ID} -> USERS {ID}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты ID, CONTENT_ID, USERS_ID, VALUE, created_at являются атомарными.
- 2 НФ: Атрибуты CONTENT_ID, USERS_ID, VALUE, created_at полностью функционально зависят от первичного ключа ID.
- 3 НФ: Атрибуты CONTENT_ID, USERS_ID, VALUE, created_at не зависят от других атрибутов.
  - НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    RATING_OF_CONTENT {
        INT ID PK "Уникальный идентификатор рейтинга контента"
        INT CONTENT_ID FK "Идентификатор контента, который имеет этот рейтинг"
        INT USERS_ID FK "Идентификатор пользователя, который оценил контент"
        INT VALUE "Рейтинг контента"
        TIME created_at "время создания кортежа, для отладки"
    }
```

## Таблица RATING_OF_PERSON

Таблица `RATING_OF_PERSON` содержит информацию о рейтингах персон.
<p> Функциональные зависимости: <p> 

- `{ID} -> {PERSON_ID, USERS_ID, VALUE, created_at}`
- `{PERSON_ID} -> PERSON {ID}`
- `{USERS_ID} -> USERS {ID}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты ID, PERSON_ID, USERS_ID, VALUE, created_at являются атомарными.
- 2 НФ: Атрибуты PERSON_ID, USERS_ID, VALUE, created_at полностью функционально зависят от первичного ключа ID.
- 3 НФ: Атрибуты PERSON_ID, USERS_ID, VALUE, created_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    RATING_OF_PERSON {
        INT ID PK "Уникальный идентификатор рейтинга персоны"
        INT PERSON_ID FK "Идентификатор персоны, которая имеет этот рейтинг"
        INT USERS_ID FK "Идентификатор пользователя, который оценил персону"
        INT VALUE "Рейтинг персоны"
        TIME created_at "время создания кортежа, для отладки"
    }
```

# Вспомогательные таблицы

## Таблица REVIEW_LIKES

Таблица `REVIEW_LIKES` содержит информацию о лайках/дизлайках комментариев. Она реализует связь между комментарием и пользователем, 
который оценил чужой комментарий.
<p> Функциональные зависимости: <p> 

- `{COMMENT_ID} -> COMMENT {ID}`
- `{USER_ID} -> USER {ID}`
- `{COMMENT_ID, USER_ID} -> {VALUE, created_at}`

```mermaid
erDiagram
    REVIEW_LIKES{
        INT REVIEW_ID FK "Идентификатор комментария"
        INT USERS_ID FK "Идентификатор пользователя, оценившего комментарий"
        INT VALUE "Значение оценки комментария (1 или -1, например)"
        TIME created_at "время создания кортежа, для отладки"
    }
```

## Таблица SAVED_PERSON

Таблица `SAVED_PERSON` содержит информацию о персонах, которые были сохранены пользователем. Она реализует связь М:N между пользователями и персонами.
<p> Функциональные зависимости: <p> 

- `{PERSON_ID} -> PERSON {ID}`
- `{USER_ID} -> USER {ID}`
- `{PERSON_ID, USER_ID} -> {created_at}`

```mermaid
erDiagram
    SAVED_PERSON{
        INT PERSON_ID FK "Идентификатор сохраненного персону"
        INT USERS_ID FK "Идентификатор пользователя, сохранившего персону"
        TIME created_at "время создания кортежа, для отладки"
    }
```

## Таблица CONTENT_STATUS

Таблица `CONTENT_STATUS` содержит информацию о контенте, которому пользователем был поставлен определенный статус.
Она реализует связь М:N между пользователями и контентом.
<p> Функциональные зависимости: <p> 

- `{CONTENT_ID} -> CONTENT {ID}`
- `{USER_ID} -> USER {ID}`
- `{STATUS_ID} -> STATUS {ID}`
- `{CONTENT_ID, USER_ID, STATUS_ID} -> {created_at}`

```mermaid
erDiagram
    CONTENT_STATUS{
        INT CONTENT_ID FK "Идентификатор контента"
        INT USERS_ID FK "Идентификатор пользователя"
        INT STATUS_ID FK "Статус контента (Viewed, Planned, Reconsidering, Favourites(избранное))"
        TIME created_at "время создания кортежа, для отладки"
    }
```

## Таблица GENRE_CONTENT

Таблица `GENRE_CONTENT` реализует связь М:N между контентом и жанрами.
<p> Функциональные зависимости: <p> 

- `{GENRE_ID} -> GENRE {ID}`
- `{CONTENT_ID} -> CONTENT {ID}`
- `{GENRE_ID, CONTENT_ID} -> {created_at}`

```mermaid
erDiagram
    GENRE_CONTENT{
        INT GENRE_ID FK "Идентификатор жанра"
        INT CONTENT_ID FK "Идентификатор контента"
        TIME created_at "время создания кортежа, для отладки"
    }
```

## Таблица COUNTRY_CONTENT

Таблица `COUNTRY_CONTENT` реализует связь М:N между контентом и странами.
<p> Функциональные зависимости: <p> 

- `{COUNTRY_ID} -> COUNTRY {ID}`
- `{CONTENT_ID} -> CONTENT {ID}`
- `{COUNTRY_ID, CONTENT_ID} -> {created_at}`

```mermaid
erDiagram
    COUNTRY_CONTENT{
        INT COUNTRY_ID FK "Идентификатор страны"
        INT CONTENT_ID FK "Идентификатор контента"
        TIME created_at "время создания кортежа, для отладки"
    }
```

## Таблица PERSON_ROLE

Таблица `PERSON_ROLE` реализует связь М:N между персонами и ролями в контенте. У одной персоны может быть несколько ролей в разных контентах.
<p> Функциональные зависимости: <p> 

- `{PERSON_ID} -> PERSON {ID}`
- `{ROLE_ID} -> ROLE {ID}`
- `{CONTENT_ID} -> CONTENT {ID}`
- `{PERSON_ID, ROLE_ID, CONTENT_ID} -> {created_at}`

```mermaid
erDiagram
    CONTENT_PERSON{
        INT CONTENT_ID FK "Идентификатор контента"
        INT PERSON_ID FK "Идентификатор персоны"
        INT ROLES_ID FK "Роль персоны в контенте"
        TIME created_at "время создания кортежа, для отладки"
    }
```