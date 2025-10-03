-- name: Create :one
INSERT INTO daily_question (
    domain,
    paragraph,
    correct,
    svg,
    prompt,
    explanation,
    difficulty,
    choice_a,
    choice_b,
    choice_c,
    choice_d,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING *;


-- name: GetAll :many
SELECT * FROM daily_question ORDER BY created_at;

-- name: GetLatest :many
SELECT *
FROM daily_question
WHERE DATE(created_at) = CURRENT_DATE
LIMIT 2;
