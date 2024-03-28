```mermaid
erDiagram
    AUDIENCE {
        INT ID PK "Уникальный идентификатор аудитории"
        INT CONTENT_ID FK "Идентификатор контента, который имеет эту аудиторию"
        INT SIZE_IN_THOUSANDS "Размер аудитории в тысячах"
        INT COUNTRY_ID "Идентификатор страны аудитории"
        TIME created_at "время создания кортежа, для отладки"
    }
    COUNTRY {
        INT ID PK "Уникальный идентификатор страны"
        STRING NAME "Название страны"
        TIME created_at "время создания кортежа, для отладки"
    }
    BOXOFFICE {
        INT ID PK "Уникальный идентификатор кассовых сборов"
        INT CONTENT_ID FK "Идентификатор контента, который имеет этот кассовый сбор"
        INT COUNTRY_ID "Идентификатор страны кассовых сборов"
        INT REVENUE "Сумма сборов"
        TIME created_at "время создания кортежа, для отладки"
    }
    GENRE {
        INT ID PK "Уникальный идентификатор жанра"
        STRING NAME "Название жанра на английском"
        STRING NAME_RU "Название жанра на русском"
        TIME created_at "время создания кортежа, для отладки"
    }
    BIRTHPLACE{
        INT ID PK "Уникальный идентификатор места рождения"
	    STRING CITY "Город рождения"
	    STRING REGION "Регион рождения"
	    INT COUNTRY_ID "Идентификатор страны рождения"
        TIME created_at "время создания кортежа, для отладки"
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
	    STRING SEX "Пол персоны"
	    STRING PHOTO "Фотография персоны"
	    INT HEIGHT "Рост персоны"         
	    STRING SPOUSE "Супруг(а) персоны"  
	    STRING CHILDREN "Дети персоны"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    ROLES{
        INT ID PK "Уникальный идентификатор роли персоны в контенте"
	    STRING NAME "Название роли (актер, режиссер, дублер и т.д.)"  
        TIME created_at "время создания кортежа, для отладки"
    }
    FILM {
        INT ID PK "Уникальный идентификатор фильма"
        INT CONTENT_ID FK "Идентификатор контента, который относится к фильму"
        INT YEAR "Год выпуска фильма"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
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
    STATUS{
        INT ID PK "Уникальный идентификатор статуса"
        STRING STATUS "Статус контента (Viewed, Planned, Reconsidering, Favourites(избранное))"
        TIME created_at "время создания кортежа, для отладки"
    }
    EPISODE {
        INT ID PK "Уникальный идентификатор эпизода"
        INT SEASON_ID FK "Идентификатор сезона, к которому относится эпизод"
        STRING DESCRIPTION "Описание эпизода"
        INT EPISODE_NUMBER "Номер эпизода"
        STRING VIEWED "Просмотрен ли эпизод"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    SEASON {
        INT ID PK "Уникальный идентификатор сезона"
        INT SERIES_ID FK "Идентификатор сериала, к которому относится сезон"
        INT YEAR_START "Год начала сезона"
        INT YEAR_END "Год окончания сезона"
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
    SERIES {
        INT ID PK "Уникальный идентификатор сериала"
        INT CONTENT_ID FK "Идентификатор контента, который относится к  сериалу"
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
        TIME created_at "время создания кортежа, для отладки"
        TIME updated_at "время изменения кортежа, для отладки"
    }
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
    NOMINATION {
        INT ID PK "Уникальный идентификатор номинации"
        STRING TITLE "Название номинации"
        INT CONTENT_ID FK "Идентификатор контента, за который дана номинация"
        INT PERSON_ID FK  "Идентификатор персоны, которой присвоена номинация"
        INT AWARD_ID FK  "Идентификатор награды"
        TIME created_at "время создания кортежа, для отладки"
    }
    AWARD {
        INT ID PK "Уникальный идентификатор награды"
        INT YEAR "Год присуждения награды"
        STRING NAME "Тип награды"
        TIME created_at "время создания кортежа, для отладки"
    }
    RATING_OF_CONTENT {
        INT ID PK "Уникальный идентификатор рейтинга контента"
        INT CONTENT_ID FK "Идентификатор контента, который имеет этот рейтинг"
        INT USERS_ID FK "Идентификатор пользователя, который оценил контент"
        INT VALUE "Рейтинг контента"
        TIME created_at "время создания кортежа, для отладки"
    }
    RATING_OF_PERSON {
        INT ID PK "Уникальный идентификатор рейтинга персоны"
        INT PERSON_ID FK "Идентификатор персоны, которая имеет этот рейтинг"
        INT USERS_ID FK "Идентификатор пользователя, который оценил персону"
        INT VALUE "Рейтинг персоны"
        TIME created_at "время создания кортежа, для отладки"
    }

    USERS ||--o{ CONTENT_STATUS : "USERS_ID"
    CONTENT ||--o{ CONTENT_STATUS : "CONTENT_ID"
    STATUS ||--o{ CONTENT_STATUS : "CONTENT_ID"
    USERS ||--o{ SAVED_PERSON : "USERS_ID"
    PERSON ||--o{ SAVED_PERSON : "PERSON_ID"
    GENRE ||--o{ GENRE_CONTENT : "GENRE_ID"
    CONTENT ||--o{ GENRE_CONTENT : "CONTENT_ID"
    COUNTRY ||--o{ COUNTRY_CONTENT : "COUNTRY_ID"
    CONTENT ||--o{ COUNTRY_CONTENT : "CONTENT_ID"
    CONTENT ||--o{ CONTENT_PERSON : "CONTENT_ID"
    PERSON ||--o{ CONTENT_PERSON : "PERSON_ID"
    ROLES ||--o{ CONTENT_PERSON : "ROLES_ID"

    CONTENT ||--o{ RATING_OF_CONTENT : "CONTENT_ID"
    USERS ||--o{ RATING_OF_CONTENT : "USERS_ID"
    PERSON ||--o{ RATING_OF_PERSON : "PERSON_ID"
    USERS ||--o{ RATING_OF_PERSON : "USERS_ID"
    BIRTHPLACE ||--o{ PERSON : "BIRTHPLACE_ID"
    SERIES ||--o{ SEASON : "SERIES_ID"
    SEASON ||--o{ EPISODE : "SEASON_ID"
    USERS ||--o{ REVIEW : "USERS_ID"
    CONTENT ||--o{ REVIEW : "CONTENT_ID"
    CONTENT ||--o{ AUDIENCE : " CONTENT_ID"
    CONTENT ||--o{ NOMINATION : "CONTENT_ID"
    AWARD ||--o{ NOMINATION : "AWARD_ID"
    PERSON ||--o{ NOMINATION : "PERSON_ID"
    CONTENT ||--o{ BOXOFFICE : "CONTENT_ID"

    REVIEW ||--o{ REVIEW_LIKES : "REVIEW_ID"
    USERS ||--o{ REVIEW_LIKES : "USERS_ID"
    
    CONTENT ||--|| SERIES : "CONTENT_ID"
    CONTENT ||--|| FILM : "CONTENT_ID"
    

    REVIEW_LIKES{
        INT REVIEW_ID FK "Идентификатор комментария"
        INT USERS_ID FK "Идентификатор пользователя, оценившего комментарий"
        INT VALUE "Значение оценки комментария (1 или -1, например)"
        TIME created_at "время создания кортежа, для отладки"
    }
    SAVED_PERSON{
        INT PERSON_ID FK "Идентификатор сохраненного персону"
        INT USERS_ID FK "Идентификатор пользователя, сохранившего персону"
        TIME created_at "время создания кортежа, для отладки"
    }
    CONTENT_STATUS{
        INT CONTENT_ID FK "Идентификатор контента"
        INT USERS_ID FK "Идентификатор пользователя"
        INT STATUS_ID FK "Статус контента (Viewed, Planned, Reconsidering, Favourites(избранное))"
        TIME created_at "время создания кортежа, для отладки"
    }
    GENRE_CONTENT{
        INT GENRE_ID FK "Идентификатор жанра"
        INT CONTENT_ID FK "Идентификатор контента"
        TIME created_at "время создания кортежа, для отладки"
    }
    COUNTRY_CONTENT{
        INT COUNTRY_ID FK "Идентификатор страны"
        INT CONTENT_ID FK "Идентификатор контента"
        TIME created_at "время создания кортежа, для отладки"
    }
    CONTENT_PERSON{
        INT CONTENT_ID FK "Идентификатор контента"
        INT PERSON_ID FK "Идентификатор персоны"
        INT ROLES_ID FK "Роль персоны в контенте"
        TIME created_at "время создания кортежа, для отладки"
    }
