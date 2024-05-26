-- +goose Up
DROP TRIGGER IF EXISTS update_content_rating ON review;
DROP FUNCTION IF EXISTS update_content_rating();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_content_rating()
    RETURNS TRIGGER AS
$$
BEGIN
    -- пересчитываем рейтинг контента по 10 балльной шкале
    -- операция затратная, но рецензии оставляют сравнительно редко, так что будем считать это допустимым
    IF TG_OP = 'DELETE' THEN
        UPDATE content
        SET rating = COALESCE((SELECT AVG(content_rating) FROM review WHERE content_id = OLD.content_id), imdb, 0)
        WHERE id = OLD.content_id;
    ELSE
        UPDATE content
        SET rating = (SELECT AVG(content_rating) FROM review WHERE content_id = NEW.content_id)
        WHERE id = NEW.content_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE 'plpgsql';
-- +goose StatementEnd

CREATE TRIGGER update_content_rating
    AFTER INSERT OR DELETE OR UPDATE
    ON review
    FOR EACH ROW
EXECUTE FUNCTION update_content_rating();

