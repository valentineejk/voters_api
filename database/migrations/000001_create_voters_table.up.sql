CREATE TABLE voters (
    id         TEXT PRIMARY KEY,
    full_name  TEXT NOT NULL,
    nin        TEXT NOT NULL UNIQUE,
    dob        DATE NOT NULL,
    state      TEXT NOT NULL,
    lga        TEXT NOT NULL,
    phone      TEXT NOT NULL,
    status     TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_voters_state ON voters (state);
CREATE INDEX idx_voters_status ON voters (status);
CREATE INDEX idx_voters_created_at ON voters (created_at DESC);
