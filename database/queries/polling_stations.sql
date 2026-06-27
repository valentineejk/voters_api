-- name: GetPollingStation :one
SELECT * FROM polling_stations
WHERE id = $1;

-- name: GetPollingStationByCode :one
SELECT * FROM polling_stations
WHERE code = $1;

-- name: ListPollingStations :many
SELECT * FROM polling_stations
WHERE (sqlc.narg('state')::text IS NULL OR state = sqlc.narg('state'))
  AND (sqlc.narg('lga')::text IS NULL OR lga = sqlc.narg('lga'))
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status'))
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountPollingStations :one
SELECT COUNT(*) FROM polling_stations
WHERE (sqlc.narg('state')::text IS NULL OR state = sqlc.narg('state'))
  AND (sqlc.narg('lga')::text IS NULL OR lga = sqlc.narg('lga'))
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status'));

-- name: CreatePollingStation :one
INSERT INTO polling_stations (id, code, name, state, lga, ward, address, latitude, longitude, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: UpdatePollingStationStatus :one
UPDATE polling_stations
SET status = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;