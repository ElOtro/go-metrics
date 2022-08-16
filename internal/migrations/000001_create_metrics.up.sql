CREATE TABLE metrics (
  id BIGSERIAL PRIMARY KEY,
  name text UNIQUE NOT NULL,
  type text NOT NULL,
  value double precision,
  delta bigint,
  hash text
);
CREATE UNIQUE INDEX metrics_name_index ON metrics USING btree (name);
CREATE INDEX metrics_type_index ON metrics USING btree (type);