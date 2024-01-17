CREATE TABLE IF NOT EXISTS url(
    id int primary key,
    alias text not null unique,
    url text not null
);

CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);