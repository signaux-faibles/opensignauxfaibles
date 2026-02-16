CREATE TABLE IF NOT EXISTS tmp_saved_indexes (
  table_name TEXT NOT NULL,
  index_name TEXT NOT NULL,
  index_def  TEXT NOT NULL,
  PRIMARY KEY (table_name, index_name)
);
---- create above / drop below ----

DROP TABLE IF EXISTS tmp_saved_indexes;
