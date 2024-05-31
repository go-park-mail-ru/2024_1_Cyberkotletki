-- +goose Up
ALTER TABLE content
ADD COLUMN streaming_url TEXT;

-- джанго освобожденный
UPDATE content SET streaming_url = 'https://kinoskop_dev.hb.ru-msk.vkcs.cloud/streaming/%D0%94%D0%B6%D0%B0%D0%BD%D0%B3%D0%BE%20%D0%BE%D1%81%D0%B2%D0%BE%D0%B1%D0%BE%D0%B6%D0%B4%D0%B5%D0%BD%D0%BD%D1%8B%D0%B9%20%D0%A4%D0%B8%D0%BB%D1%8C%D0%BC.mp4' WHERE id = 19;
-- мстители
UPDATE content SET streaming_url = 'https://kinoskop_dev.hb.ru-msk.vkcs.cloud/streaming/%D0%9C%D1%81%D1%82%D0%B8%D1%82%D0%B5%D0%BB%D0%B8%20%D0%A4%D0%B8%D0%BB%D1%8C%D0%BC.mp4' WHERE id = 615;
-- серый человек
UPDATE content SET streaming_url = 'https://kinoskop_dev.hb.ru-msk.vkcs.cloud/streaming/%D0%A1%D0%B5%D1%80%D1%8B%D0%B9%20%D1%87%D0%B5%D0%BB%D0%BE%D0%B2%D0%B5%D0%BA%20%D0%A4%D0%B8%D0%BB%D1%8C%D0%BC.mp4' WHERE id = 1784;
-- шрек
UPDATE content SET streaming_url = 'https://kinoskop_dev.hb.ru-msk.vkcs.cloud/streaming/%D0%A8%D1%80%D0%B5%D0%BA%201%20%D0%A4%D0%B8%D0%BB%D1%8C%D0%BC.mp4' WHERE id = 14;
-- шрек 2
UPDATE content SET streaming_url = 'https://kinoskop_dev.hb.ru-msk.vkcs.cloud/streaming/%D0%A8%D1%80%D0%B5%D0%BA%202%20%D0%A4%D0%B8%D0%BB%D1%8C%D0%BC.mp4' WHERE id = 1299;
