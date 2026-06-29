DO $$
BEGIN
  CREATE EXTENSION IF NOT EXISTS pg_parquet;
EXCEPTION WHEN OTHERS THEN
  RAISE NOTICE 'pg_parquet extension not available, skipping: %', SQLERRM;
END $$;

---- create above / drop below ----

DROP EXTENSION IF EXISTS pg_parquet;
