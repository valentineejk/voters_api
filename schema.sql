CREATE SEQUENCE IF NOT EXISTS voter_id_seq;

CREATE TABLE voters (
    id         TEXT PRIMARY KEY,
    voter_id   TEXT NOT NULL UNIQUE DEFAULT ('VTR-' || lpad(nextval('voter_id_seq')::text, 5, '0')),
    full_name  TEXT NOT NULL,
    nin        TEXT NOT NULL UNIQUE,
    dob        DATE NOT NULL,
    state      TEXT NOT NULL,
    lga        TEXT NOT NULL,
    phone      TEXT NOT NULL,
    status     TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE polling_stations (
    id         TEXT PRIMARY KEY,                     
    code       TEXT NOT NULL UNIQUE,                 
    name       TEXT NOT NULL,
    state      TEXT NOT NULL,
    lga        TEXT NOT NULL,
    ward       TEXT NOT NULL,
    address    TEXT,
    latitude   NUMERIC(9,6),
    longitude  NUMERIC(9,6),
    status     TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


ALTER TABLE voters ADD COLUMN polling_station_id TEXT REFERENCES polling_stations(id) ON DELETE SET NULL;