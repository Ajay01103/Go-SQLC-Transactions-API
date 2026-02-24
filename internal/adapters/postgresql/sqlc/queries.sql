-- name: GetUserByEmail :one
SELECT id, name, email, password, profile_picture, created_at, updated_at
FROM users
WHERE email = $1
LIMIT 1;

-- name: GetUserByID :one
SELECT id, name, email, password, profile_picture, created_at, updated_at
FROM users
WHERE id = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
    id,
    name,
    email,
    password,
    profile_picture
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;