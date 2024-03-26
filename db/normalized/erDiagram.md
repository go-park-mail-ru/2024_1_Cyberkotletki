```mermaid
erDiagram
    AUDIENCE {
        INT ID PK "Уникальный идентификатор аудитории"
        INT CONTENT_ID FK "Идентификатор контента, который имеет эту аудиторию"
        INT SIZE_IN_THOUSANDS "Размер аудитории в тысячах"
        STRING COUNTRY "страна аудитории"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    COUNTRY {
        INT ID PK "Уникальный идентификатор страны"
        STRING NAME "Название страны"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    BOXOFFICE {
        INT ID PK "Уникальный идентификатор кассовых сборов"
        INT CONTENT_ID FK "Идентификатор контента, который имеет этот кассовый сбор"
        STRING COUNTRY "Страна кассовых сборов"
        INT REVENUE "Сумма сборов"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    GENRE {
        INT ID PK "Уникальный идентификатор жанра"
        STRING NAME "Название жанра на английском"
        STRING NAME_RU "Название жанра на русском"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    BIRTHPLACE{
        INT ID PK "Уникальный идентификатор места рождения"
	    STRING CITY "Город рождения"
	    STRING REGION "Регион рождения"
	    STRING COUNTRY "страна аудитории"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    PERSON {
        INT ID PK "Уникальный идентификатор персоны"
        STRING FIRST_NAME "Имя персоны"  
	    STRING LAST_NAME "Фамилия персоны"  
	    TIME BIRTH_DATE "Дата рождения персоны"  
        INT BIRTHPLACE_ID FK "Идентификатор места рождения персоны"  
	    TIME DEATH_DATE "Дата смерти персоны"  
	    TIME START_CAREER "Дата начала карьеры персоны" 
	    TIME END_CAREER "Дата окончания карьеры персоны" 
	    STRING PHOTO "Фотография персоны"  
        STRING GENDER "Пол - М/Ж" 
	    INT HEIGHT "Рост персоны"         
	    STRING SPOUSE "Супруг(а) персоны"  
	    STRING CHILDREN "Дети персоны"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    ROLE{
        INT ID PK "Уникальный идентификатор роли персоны в контенте"
	    STRING NAME "Название роли (актер, режиссер, дублер и т.д.)"  
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    FILM {
        INT ID PK "Уникальный идентификатор фильма"
        INT CONTENT_ID FK "Идентификатор контента, который относится к фильму"
        INT YEAR "Год выпуска фильма"
        INT DURATION "Продолжительность фильма в минутах"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    CONTENT {
        INT ID PK "Уникальный идентификатор контента"
        STRING TITLE "Название контента"
        STRING ORIGINAL_TITLE "Оригинальное название контента"
        INT BUDGET "Бюджет контента"
        INT MARKETING "Маркетинговые затраты на контент"
        TIME PREMIERE "Дата премьеры контента"
        TIME RELEASE "Дата выпуска контента"
        INT AGE_RESTRICTION "Возрастное ограничение контента"
        INT IMDB "Рейтинг IMDB контента"
        STRING DESCRIPTION "Описание контента"
        STRING POSTER "Постер контента"
        STRING PLAYBACK "Воспроизведение на заднем плане небольшого фрагмента видео контента"
        STRING TYPE "F - Film, S - Season"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    STATUS{
        INT ID PK "Уникальный идентификатор статуса"
        STRING STATUS "Статус контента (Viewed, Planned, Reconsidering, Favourites(избранное))"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    EPISODE {
        INT ID PK "Уникальный идентификатор эпизода"
        INT SEASON_ID FK "Идентификатор сезона, к которому относится эпизод"
        STRING DESCRIPTION "Описание эпизода"
        INT EPISODE_NUMBER "Номер эпизода"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    SEASON {
        INT ID PK "Уникальный идентификатор сезона"
        INT CONTENT_ID FK "Идентификатор контента, который относится к  сезону"
        INT SERIES_ID FK "Идентификатор сериала, к которому относится сезон"
        INT YEAR_START "Год начала сезона"
        INT YEAR_END "Год окончания сезона"
        INT COUNT_EPISODES "Число эпизодов в сезоне (выпущено или запланировано)"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    SERIES {
        INT ID PK "Уникальный идентификатор сериала"
        STRING TITLE "Название сериала"
        INT YEAR_START "Год начала сериала"
        INT YEAR_END "Год окончания сериала"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    USERS {
        INT ID PK "Уникальный идентификатор пользователя"
        STRING NAME "Имя пользователя"
        STRING EMAIL "Электронная почта пользователя"
        STRING PASSWORD_HASHED "Хэш пароля пользователя"
        STRING SALT_PASSWORD "Соль для генерации хэша пароля"
        TIME BIRTH_DATE "День рождения пользователя"
        TIME DATE_REGISTERED "Дата регистрации пользователя"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    COMMENT {
        INT ID PK "Уникальный идентификатор комментария"
        INT USER_ID FK "Идентификатор пользователя, оставившего комментарий"
        INT CONTENT_ID FK "Идентификатор контента, к которому относится комментарий"
        STRING TITLE "Заголовок комментария"
        STRING TEXT "Текст комментария"
        INT RATING_USER "Оценка контента пользователем, оставившим комментарий"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    NOMINATION {
        INT ID PK "Уникальный идентификатор номинации"
        STRING TITLE "Название номинации"
        INT CONTENT_ID FK "Идентификатор контента, за который дана номинация"
        INT PERSON_ID FK  "Идентификатор персоны, которой присвоена номинация"
        INT AWARD_ID FK  "Идентификатор награды"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    AWARD {
        INT ID PK "Уникальный идентификатор награды"
        INT YEAR "Год присуждения награды"
        STRING NAME "Тип награды"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    RATING {
        INT ID PK  "Уникальный идентификатор рейтинга"
        INT USER_ID FK "Идентификатор пользователя, оставившего рейтинг"
        INT VALUE "Значение рейтинга"
        STRING TYPE "C - content, P - person"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    RATING_TYPE{
        INT RATING_ID FK  "Идентификатор рейтинга"
        INT ENTITY_ID FK  "Идентификатор контента или персоны, которому присвоен рейтинг"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }

    USER ||--o{ CONTENT_STATUS : "USER_ID"
    CONTENT ||--o{ CONTENT_STATUS : "CONTENT_ID"
    STATUS ||--o{ CONTENT_STATUS : "CONTENT_ID"
    USER ||--o{ SAVED_PERSON : "USER_ID"
    PERSON ||--o{ SAVED_PERSON : "PERSON_ID"
    GENRE ||--o{ GENRE_CONTENT : "GENRE_ID"
    CONTENT ||--o{ GENRE_CONTENT : "CONTENT_ID"
    COUNTRY ||--o{ COUNTRY_CONTENT : "COUNTRY_ID"
    CONTENT ||--o{ COUNTRY_CONTENT : "CONTENT_ID"
    CONTENT ||--o{ CONTENT_PERSON : "CONTENT_ID"
    PERSON ||--o{ CONTENT_PERSON : "PERSON_ID"
    ROLE ||--o{ CONTENT_PERSON : "ROLE_ID"

    CONTENT ||--o{ RATING_TYPE : "ENTITY_ID=CONTENT_ID"
    PERSON ||--o{ RATING_TYPE : "ENTITY_ID=PERSON_ID"
    BIRTHPLACE ||--o{ PERSON : "BIRTHPLACE_ID"
    SERIES ||--o{ SEASON : "SERIES_ID"
    SEASON ||--o{ EPISODE : "SEASON_ID"
    USER ||--o{ COMMENT : "USER_ID"
    CONTENT ||--o{ COMMENT : "CONTENT_ID"
    USER ||--o{ RATING : "USER_ID"
    CONTENT ||--o{ AUDIENCE : " CONTENT_ID"
    CONTENT ||--o{ NOMINATION : "CONTENT_ID"
    AWARD ||--o{ NOMINATION : "AWARD_ID"
    PERSON ||--o{ NOMINATION : "PERSON_ID"
    CONTENT ||--o{ BOXOFFICE : "CONTENT_ID"

    COMMENT ||--o{ COMMENT_LIKES : "COMMENT_ID"
    USER ||--o{ COMMENT_LIKES : "USER_ID"
    
    CONTENT ||--|| FILM : "CONTENT_ID"
    CONTENT ||--|| SEASON : "CONTENT_ID"
    

    COMMENT_LIKES{
        INT COMMENT_ID FK "Идентификатор комментария"
        INT USER_ID FK "Идентификатор пользователя, оценившего комментарий"
        INT VALUE "Значение оценки комментария (1 или -1, например)"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    SAVED_PERSON{
        INT PERSON_ID FK "Идентификатор сохраненного персону"
        INT USER_ID FK "Идентификатор пользователя, сохранившего персону"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    CONTENT_STATUS{
        INT CONTENT_ID FK "Идентификатор контента"
        INT USER_ID FK "Идентификатор пользователя"
        INT STATUS_ID FK "Статус контента (Viewed, Planned, Reconsidering, Favourites(избранное))"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    GENRE_CONTENT{
        INT GENRE_ID FK "Идентификатор жанра"
        INT CONTENT_ID FK "Идентификатор контента"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    COUNTRY_CONTENT{
        INT COUNTRY_ID FK "Идентификатор страны"
        INT CONTENT_ID FK "Идентификатор контента"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    CONTENT_PERSON{
        INT CONTENT_ID FK "Идентификатор контента"
        INT PERSON_ID FK "Идентификатор персоны"
        INT ROLE_ID FK "Роль персоны в контенте"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
