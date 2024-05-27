-- +goose Up
-- Удаление таблиц
DROP TABLE IF EXISTS genre_ongoing_content;
DROP TABLE IF EXISTS ongoing_content;

-- Добавление нового поля в таблицу content
ALTER TABLE content ADD COLUMN ongoing BOOLEAN DEFAULT false NOT NULL;
ALTER TABLE content ADD COLUMN ongoing_date TIMESTAMPTZ DEFAULT NULL;
ALTER TABLE content ADD CONSTRAINT check_ongoing_date CHECK ((ongoing = false) OR (ongoing = true AND ongoing_date IS NOT NULL));

-- Теперь может быть NULL
ALTER TABLE content DROP CONSTRAINT IF EXISTS content_poster_upload_id_fkey;
ALTER TABLE content ADD CONSTRAINT content_poster_upload_id_fkey FOREIGN KEY (poster_upload_id) REFERENCES static (id) ON DELETE SET NULL;

-- Static
INSERT INTO static (name, path) VALUES
    ('11.webp', 'ongoing')
;

UPDATE content SET ongoing = true, ongoing_date = MAKE_TIMESTAMPTZ(2024, 6, 6, 0, 0, 0) WHERE kinopoisk_id = 1236605;
UPDATE content SET ongoing = true, ongoing_date = MAKE_TIMESTAMPTZ(2024, 6, 6, 0, 0, 0) WHERE kinopoisk_id = 5024418;

-- Content
INSERT INTO content (kinopoisk_id, content_type, title, original_title, slogan, age_restriction, imdb, rating, description, poster_upload_id, trailer_url, ongoing_date, ongoing)
VALUES
    (1280694, 'movie', 'Пришельцы', 'Oegye+in 1bu', 'Отправиться в прошлое, чтобы спасти будущее', 12, 0, 0, 'На протяжении многих веков инопланетная раса содержит своих преступников в человеческих телах, а люди об этом даже не догадываются. Способные перемещаться во времени с помощью энергетического кинжала Охранник с роботом-помощником поставлены следить, чтобы заключённые не пришли в себя и не сбежали. В 1380 году они ловят очередного проснувшегося преступника, но у его человеческой оболочки остаётся младенец, очевидно обречённый на смерть. 1391 год, Корё. Обладающий волшебным веером и навыками боевых искусств даос Му-рык в компании двух помощников ищет волшебный клинок, за который назначена немалая награда, и в процессе сталкивается с загадочной девушкой. А в 2022 году в Сеуле 10-летняя «дочь» Охранника уже давно подозревает, что «папа» — не человек, и начинает вынюхивать, чем же он таким интересным занят.', (SELECT id FROM static WHERE path='ongoing' AND name='3.webp'), 'https://www.youtube.com/embed/3waDpx6aU1I', MAKE_TIMESTAMPTZ(2024, 6, 20, 0, 0, 0), true),
    (660906, 'movie', 'Самая большая луна', 'Самая большая луна', NULL, 6, 0, 0, 'Среди нас живут эмеры, сверхлюди, наделённые необычными способностями. Они могут управлять человеческими эмоциями и не чувствуют боли, но при этом не способны любить. Их существование находится под угрозой, и Денис, молодой эмер-полукровка, отправляется на поиски избранной, которая может всё изменить. Катю всю жизнь ограждали от любых эмоций, она даже не представляет, какой невероятной силой обладает.', (SELECT id FROM static WHERE path='ongoing' AND name='4.webp'), 'https://www.youtube.com/embed/1tpmenl6t38', MAKE_TIMESTAMPTZ(2024, 7, 4, 0, 0, 0), true),
    (919272, 'movie', 'Крысолов. Древнее проклятие', 'Sonnim', NULL, 12, 0, 0, '1950-е. Отец с больным сыном по дороге в Сеул решают заглянуть в отдаленную деревню. Вскоре они узнают, что деревню наводнили огромные крысы, которые питаются человечиной, а местный староста скрывает от жителей страшную тайну. Покидать деревню никому не разрешается.', (SELECT id FROM static WHERE path='ongoing' AND name='5.webp'), NULL, MAKE_TIMESTAMPTZ(2024, 7, 11, 0, 0, 0), true),
    (968031, 'movie', 'Бордерлендс', 'Borderlands', 'Хаос любит компанию', 18, 0, 0, 'Лилит, одна из самых известных охотниц за головами во Вселенной, терпеть не может две вещи: Пандору — свою родную планету, и Атласа — редкого отброса и влиятельного бандита. И все же ей придется иметь дело и с тем, и с другим: отправившись на поиски дочери Атласа на Пандору, она объединяется с другими искателями приключений. Вместе эти далеко не герои должны сразиться с инопланетной расой, жестокими охотниками, раскрыть один из самых страшных секретов планеты и, в придачу, спасти мир от надвигающегося зла. Если кто и способен все это провернуть со стилем, то это они.', (SELECT id FROM static WHERE path='ongoing' AND name='6.webp'), 'https://www.youtube.com/embed/tDv0WYvxZgM', MAKE_TIMESTAMPTZ(2024, 8, 8, 0, 0, 0), true),
    (441406, 'movie', 'Ворон', 'The Crow', 'A modern reimagining of the beloved character, THE CROW, based on the original graphic novel by James O''Barr', 18, 0, 0, 'Пожертвовав собой, чтобы спасти возлюбленную, Эрик Дрэйвен застревает между мирами живых и мёртвых. Вернувшись с того света, он не остановится ни перед чем, чтобы свести счёты с убийцами. Отныне он Ворон, жаждущий справедливости, и его месть будет жестока как никогда.', (SELECT id FROM static WHERE path='ongoing' AND name='7.webp'), 'https://www.youtube.com/embed/6sx1Chliz48', MAKE_TIMESTAMPTZ(2024, 8, 21, 0, 0, 0), true),
    (5047471, 'movie', 'Волшебник Изумрудного города', 'Волшебник Изумрудного города', NULL, 6, 0, 0, 'В далёком городе живёт девочка Элли. Однажды злая колдунья Гингема наколдовала ураган, который унёс Элли и ее собачку Тотошку в страну Жевунов. Чтобы вернуться домой, Элли вместе с друзьями — Страшилой, Железным Дровосеком и Трусливым Львом — отправится по желтой кирпичной дороге в Изумрудный город на поиски Волшебника, который исполнит их заветные желания.', (SELECT id FROM static WHERE path='ongoing' AND name='8.webp'), 'https://www.youtube.com/embed/p3-d9qg9Ew4', MAKE_TIMESTAMPTZ(2025, 1, 1, 0, 0, 0), true),
    (5230098, 'movie', 'Финист. Первый богатырь', 'Финист. Первый богатырь', NULL, 6, 0, 0, 'Финист Ясный Сокол — самый удалой богатырь Белогорья, самый сильный, самый ловкий и самый красивый. Все остальные богатыри на него равняются, дети хотят быть на него похожими, а девушки — просто заглядываются.', (SELECT id FROM static WHERE path='ongoing' AND name='9.webp'), 'https://www.youtube.com/embed/d_QR4hWWKzo', MAKE_TIMESTAMPTZ(2025, 1, 1, 0, 0, 0), true),
    (5304486, 'movie', 'Ждун в кино', 'Ждун в кино', NULL, 6, 0, 0, 'Корабль доброго инопланетянина Ждуна потерпел крушение во время метеоритного дождя и упал в лесу, недалеко от Абрау-Дюрсо. Ждун посылает сигнал бедствия в космос и начинает ждать помощи, но есть незадача — до его планеты 5 световых лет и ждать пришлось бы долго. К счастью, на помощь Ждуну приходит любознательный мальчик Никита, а позже подтягивается и вся его семья. Никите и семье Семеновых предстоит помочь пришельцу отремонтировать корабль, избежать козней местного афериста-бизнесмена, который страстно мечтает завладеть инопланетными технологиями, и вернуться домой. А Ждун заново научит семью Семеновых слышать и понимать друг друга, и покажет, что в чрезвычайной ситуации даже обычный человек может стать супергероем.', (SELECT id FROM static WHERE path='ongoing' AND name='10.webp'), 'https://www.youtube.com/embed/0y52J3R7J3Q', MAKE_TIMESTAMPTZ(2025, 1, 1, 0, 0, 0), true),
    (1358438, 'movie', 'Сказочные выходные', 'Сказочные выходные', NULL, 6, 0, 0, 'Обычная российская семья с обычными проблемами отправляется в диджитал-детокс отель в Муромской глуши в день, когда полчища сказочных чудищ вдруг начинают материализовываться в реальном мире.', (SELECT id FROM static WHERE path='ongoing' AND name='11.webp'), NULL, MAKE_TIMESTAMPTZ(2025, 7, 7, 0, 0, 0), true)
;


-- Countries
INSERT INTO country_content (country_id, content_id) VALUES
    (14, (SELECT id FROM content WHERE kinopoisk_id = 1280694)), -- Южная Корея
    (22, (SELECT id FROM content WHERE kinopoisk_id = 660906)), -- Россия
    (14, (SELECT id FROM content WHERE kinopoisk_id = 919272)), -- Южная Корея
    (2, (SELECT id FROM content WHERE kinopoisk_id = 968031)), -- США
    (2, (SELECT id FROM content WHERE kinopoisk_id = 441406)), -- США
    (22, (SELECT id FROM content WHERE kinopoisk_id = 5047471)), -- Россия
    (22, (SELECT id FROM content WHERE kinopoisk_id = 5230098)), -- Россия
    (22, (SELECT id FROM content WHERE kinopoisk_id = 5304486)), -- Россия
    (22, (SELECT id FROM content WHERE kinopoisk_id = 1358438)); -- Россия

-- Movies
INSERT INTO movie (content_id) VALUES
    ((SELECT id FROM content WHERE kinopoisk_id = 1280694)),
    ((SELECT id FROM content WHERE kinopoisk_id = 660906)),
    ((SELECT id FROM content WHERE kinopoisk_id = 919272)),
    ((SELECT id FROM content WHERE kinopoisk_id = 968031)),
    ((SELECT id FROM content WHERE kinopoisk_id = 441406)),
    ((SELECT id FROM content WHERE kinopoisk_id = 5047471)),
    ((SELECT id FROM content WHERE kinopoisk_id = 5230098)),
    ((SELECT id FROM content WHERE kinopoisk_id = 5304486)),
    ((SELECT id FROM content WHERE kinopoisk_id = 1358438));

-- Genres
INSERT INTO genre_content (genre_id, content_id) VALUES (5, (SELECT id FROM content WHERE kinopoisk_id = 1280694));
INSERT INTO genre_content (genre_id, content_id) VALUES (6, (SELECT id FROM content WHERE kinopoisk_id = 1280694));
--
INSERT INTO genre_content (genre_id, content_id) VALUES (8, (SELECT id FROM content WHERE kinopoisk_id = 660906));
INSERT INTO genre_content (genre_id, content_id) VALUES (7, (SELECT id FROM content WHERE kinopoisk_id = 660906));
--
INSERT INTO genre_content (genre_id, content_id) VALUES (10, (SELECT id FROM content WHERE kinopoisk_id = 919272));
INSERT INTO genre_content (genre_id, content_id) VALUES (19, (SELECT id FROM content WHERE kinopoisk_id = 919272));
--
INSERT INTO genre_content (genre_id, content_id) VALUES (5, (SELECT id FROM content WHERE kinopoisk_id = 968031));
INSERT INTO genre_content (genre_id, content_id) VALUES (6, (SELECT id FROM content WHERE kinopoisk_id = 968031));
--
INSERT INTO genre_content (genre_id, content_id) VALUES (6, (SELECT id FROM content WHERE kinopoisk_id = 441406));
INSERT INTO genre_content (genre_id, content_id) VALUES (10, (SELECT id FROM content WHERE kinopoisk_id = 441406));
--
INSERT INTO genre_content (genre_id, content_id) VALUES (7, (SELECT id FROM content WHERE kinopoisk_id = 5047471));
INSERT INTO genre_content (genre_id, content_id) VALUES (8, (SELECT id FROM content WHERE kinopoisk_id = 5047471));
--
INSERT INTO genre_content (genre_id, content_id) VALUES (7, (SELECT id FROM content WHERE kinopoisk_id = 5230098));
INSERT INTO genre_content (genre_id, content_id) VALUES (8, (SELECT id FROM content WHERE kinopoisk_id = 5230098));
--
INSERT INTO genre_content (genre_id, content_id) VALUES (2, (SELECT id FROM content WHERE kinopoisk_id = 5304486));
INSERT INTO genre_content (genre_id, content_id) VALUES (5, (SELECT id FROM content WHERE kinopoisk_id = 5304486));
--
INSERT INTO genre_content (genre_id, content_id) VALUES (9, (SELECT id FROM content WHERE kinopoisk_id = 1358438));

-- Подписка на уведомление о выходе контента
CREATE TABLE IF NOT EXISTS ongoing_subscribe (
    id  INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    user_id INT NOT NULL,
    content_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES content (id) ON DELETE CASCADE,
    UNIQUE (user_id, content_id)
);

-- fix для названия
UPDATE compilation SET title = 'Величайшие сериалы XXI века' WHERE id = 27;
UPDATE compilation SET title = 'Лучшие сериалы' WHERE id = 29;




