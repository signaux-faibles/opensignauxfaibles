-- La migration 046 enveloppait CREATE EXTENSION dans un bloc DO/EXCEPTION qui
-- avalait silencieusement les erreurs, marquant la migration appliquée même
-- quand pg_parquet n'était pas chargée. On re-tente proprement ici.

CREATE EXTENSION IF NOT EXISTS pg_parquet;

---- create above / drop below ----

DROP EXTENSION IF EXISTS pg_parquet;
