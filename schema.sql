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