-- name: Create :one
WITH new_practice AS (
    INSERT INTO practice (title)
    VALUES ($1)
    RETURNING id, title, created_at, status
),
inserted_modules AS (
    INSERT INTO module (practice_id, name)
    SELECT id, name
    FROM new_practice, (VALUES 
        ('Reading And Writing 1'),
        ('Reading And Writing 2'),
        ('Math 1'),
        ('Math 2')
    ) AS t(name)
    RETURNING id
)
SELECT * FROM new_practice;



-- name: GetPracticePreviews :many
SELECT * FROM practice;

-- name: GetById :one
SELECT id, title, created_at
FROM practice 
WHERE id = $1;


-- name: Delete :one
DELETE FROM practice WHERE id = $1 RETURNING *;
