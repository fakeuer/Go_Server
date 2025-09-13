-- name: GetChirpById :one
SELECT * FROM chirps
WHERE id = $1;

-- name: GetChirpID_ByUserID :one
SELECT * FROM chirps
WHERE user_id = $1 and id = $2;

-- name: DeleteChirpById :one
DELETE FROM chirps
WHERE id = $1
RETURNING id;
