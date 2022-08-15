DROP TABLE IF EXISTS metrics;
CREATE TABLE metrics (
  id BIGSERIAL PRIMARY KEY,
  name text NOT NULL,
  type text NOT NULL,
  delta double precision DEFAULT 0.0000,
  value bigint DEFAULT 0,
  hash text
);
CREATE INDEX metrics_name_index ON metrics USING btree (name);
CREATE INDEX metrics_type_index ON metrics USING btree (type);