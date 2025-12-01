ALTER TABLE import_logs ADD COLUMN end_date TIMESTAMP;
---- create above / drop below ----

ALTER TABLE import_logs DROP COLUMN end_date;
