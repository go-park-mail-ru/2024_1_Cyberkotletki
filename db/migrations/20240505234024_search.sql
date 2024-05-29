-- +goose Up
CREATE EXTENSION IF NOT EXISTS  pg_trgm;

-- скажем пока-пока памяти на диске
CREATE INDEX IF NOT EXISTS trgm_content_title_idx ON content USING gist (title gist_trgm_ops);
CREATE INDEX IF NOT EXISTS trgm_content_original_title_idx ON content USING gist (original_title gist_trgm_ops);
CREATE INDEX IF NOT EXISTS trgm_person_name_idx ON person USING gist (name gist_trgm_ops);
CREATE INDEX IF NOT EXISTS trgm_person_en_name_idx ON person USING gist (en_name gist_trgm_ops);
