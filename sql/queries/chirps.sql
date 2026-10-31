-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    Now(),
    Now(),
    $1,
    $2
)
RETURNING *;

-- name: GetChirp :one
SELECT * FROM chirps
WHERE id = $1;

-- name: ListChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: ListChirpsByAuthor :many
SELECT * FROM chirps
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;
