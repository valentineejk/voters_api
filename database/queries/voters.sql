-- name: GetVoter :one
SELECT * FROM voters
WHERE id = $1;

-- name: GetVoterByNIN :one
SELECT * FROM voters
WHERE nin = $1;

-- name: ExistsVoterByNIN :one
SELECT EXISTS (
    SELECT 1 FROM voters WHERE nin = $1
) AS exists;

-- name: ListVoters :many
SELECT * FROM voters
WHERE (sqlc.narg('state')::text IS NULL OR state = sqlc.narg('state'))
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status'))
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountVoters :one
SELECT COUNT(*) FROM voters
WHERE (sqlc.narg('state')::text IS NULL OR state = sqlc.narg('state'))
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status'));

-- name: CreateVoter :one
INSERT INTO voters (id, full_name, nin, dob, state, lga, phone)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateVoterStatus :one
UPDATE voters
SET status = $2
WHERE id = $1
RETURNING *;

-- name: DeleteVoter :exec
DELETE FROM voters
WHERE id = $1;
