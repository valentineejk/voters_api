CREATE SEQUENCE IF NOT EXISTS voter_id_seq;

ALTER TABLE voters
    ADD COLUMN voter_id TEXT NOT NULL UNIQUE
    DEFAULT ('VTR-' || lpad(nextval('voter_id_seq')::text, 5, '0'));
