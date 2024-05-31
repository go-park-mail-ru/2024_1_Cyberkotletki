-- +goose Up
CREATE INDEX idx_genre_content_genre_id ON genre_content (genre_id);
CREATE INDEX idx_genre_content_content_id ON genre_content (content_id);

CREATE INDEX idx_country_content_country_id ON country_content (country_id);
CREATE INDEX idx_country_content_content_id ON country_content (content_id);

CREATE INDEX idx_compilation_content_compilation_id ON compilation_content (compilation_id);
CREATE INDEX idx_compilation_content_content_id ON compilation_content (content_id);

CREATE INDEX idx_person_role_role_id ON person_role (role_id);
CREATE INDEX idx_person_role_person_id ON person_role (person_id);
CREATE INDEX idx_person_role_content_id ON person_role (content_id);

CREATE INDEX idx_review_vote_review_id ON review_vote (review_id);
CREATE INDEX idx_review_vote_user_id ON review_vote (user_id);
