CREATE TABLE events(
  ID SERIAL,
  aggregate_id uuid,
  aggregate_type VARCHAR NOT NULL,
  version int,
  type VARCHAR NOT NULL,
  data JSON
)