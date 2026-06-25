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

ALTER TABLE voters 
    ADD COLUMN polling_station_id TEXT REFERENCES polling_stations(id) ON DELETE SET NULL;

CREATE INDEX idx_polling_stations_state_lga ON polling_stations (state, lga);
CREATE INDEX idx_polling_stations_status ON polling_stations (status);
CREATE INDEX idx_voters_polling_station_id ON voters (polling_station_id) WHERE polling_station_id IS NOT NULL;