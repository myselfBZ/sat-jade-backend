-- name: GetUserById :one
SELECT  * FROM users WHERE id = $1;

-- name: Create :one
INSERT INTO users(email, role, full_name, password_hash) VALUES($1, $2, $3, $4) RETURNING *;

-- name: GetByEmail :one
SELECT * FROM users WHERE email = $1;