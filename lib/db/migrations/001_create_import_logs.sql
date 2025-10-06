CREATE TABLE IF NOT EXISTS import_logs (
  start_date     TIMESTAMP,
  parser         TEXT,
  batch_key      TEXT,
  head_fatal     TEXT[],
  head_rejected  TEXT[],
  is_fatal       BOOLEAN,
  lines_parsed   INT,
  lines_rejected INT,
  lines_skipped  INT,
  lines_valid    INT,
  summary        TEXT
)
---- create above / drop below ----

DROP TABLE import_logs;
