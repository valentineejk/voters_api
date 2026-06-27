ALTER TABLE voters DROP COLUMN IF EXISTS polling_station_id;

DROP INDEX IF EXISTS idx_voters_polling_station_id;
DROP INDEX IF EXISTS idx_polling_stations_state_lga;
DROP INDEX IF EXISTS idx_polling_stations_status;

DROP TABLE IF EXISTS polling_stations;