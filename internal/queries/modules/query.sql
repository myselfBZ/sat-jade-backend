-- name: GetByPracticeId :many
SELECT *
FROM module 
WHERE practice_id = $1
ORDER BY id;

-- name: GetByID :one
SELECT * from module WHERE id = $1;

-- name: GetByNameAndPracticeID :one 
SELECT * from module WHERE practice_id = $1 AND name = $2;