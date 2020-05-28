CREATE TABLE IF NOT EXISTS spec (
    id TEXT NOT NULL,
    column_name TEXT NOT NULL,
    data_type TEXT NOT NULL,
    PRIMARY KEY(id, column_name)
);

CREATE TABLE IF NOT EXISTS data (
    id BIGSERIAL PRIMARY KEY, 
    spec_id TEXT NOT NULL,
    content JSONB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_data_spec_id on data(spec_id);


-- init spec
INSERT INTO spec VALUES('1', 'name', 'TEXT') ON CONFLICT DO NOTHING;
INSERT INTO spec VALUES('1', 'valid', 'BOOLEAN') ON CONFLICT DO NOTHING;
INSERT INTO spec VALUES('1', 'count', 'INTEGER') ON CONFLICT DO NOTHING;

