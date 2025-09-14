-- name: GetAllChirps :many
SELECT * FROM chirps 
ORDER BY created_at ASC;

-- name: GetAllChirpsDesc :many
SELECT * FROM chirps
ORDER BY created_at DESC;

-- name: GetChirpsByAuthorIDDesc :many
SELECT * FROM chirps
WHERE user_id = $1
ORDER BY created_at DESC;
