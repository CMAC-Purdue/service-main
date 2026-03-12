-- name: GetOfficer :one
SELECT * FROM officers WHERE id = $1;

-- name: ListOfficers :many
SELECT * FROM officers ORDER BY id;

-- name: CreateOfficer :one
INSERT INTO officers (name, title, linkedin, image_uri)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteOfficer :exec
DELETE FROM officers WHERE id = $1;

