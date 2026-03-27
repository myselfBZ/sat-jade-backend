-- name: Create :one
INSERT INTO results (
    practice_id, 
    user_id, 
    english_score, 
    math_score, 
    total_score
) VALUES ($1, $2, $3, $4, $5) 
RETURNING *;

-- name: GetByID :one
SELECT 
    id, 
    user_id, 
    practice_id, 
    created_at, 
    english_score, 
    math_score, 
    total_score 
FROM results 
WHERE id = $1;

-- name: GetAll :many
SELECT 
    id, 
    user_id, 
    practice_id, 
    created_at, 
    english_score, 
    math_score, 
    total_score 
FROM results
ORDER BY created_at DESC;


-- name: GetByUserID :many
SELECT 
    r.id,
    p.title AS practice_title,
    COUNT(CASE WHEN ra.status = 'correct' THEN 1 END) AS correct_answers,
    COUNT(ra.id) AS total_questions,
    r.created_at,
    r.english_score,
    r.math_score,
    r.total_score
FROM results r
JOIN practice p ON r.practice_id = p.id
LEFT JOIN result_answers ra ON r.id = ra.result_id
WHERE r.user_id = $1
GROUP BY r.id, p.title, r.created_at
ORDER BY r.created_at DESC;

-- name: Delete :one
DELETE FROM results 
WHERE id = $1 AND user_id = $2 
RETURNING *;
