-- name: Create :one
INSERT INTO test_session(
    practice_id, 
    user_id, 
    english_score, 
    math_score, 
    total_score
    ) VALUES($1, $2, $3, $4, $5) RETURNING *;


-- name: GetByID :one
SELECT * FROM test_session WHERE id = $1;


-- name: GetAll :many
SELECT * from test_session;

-- name: GetByUserID :many
SELECT 
    ts.id,
    p.title AS practice_title,
    COUNT(CASE WHEN tsa.status = 'correct' THEN 1 END) AS correct_answers,
    COUNT(tsa.id) AS total_questions,
    ts.created_at,
    ts.english_score,
    ts.math_score,
    ts.total_score
FROM test_session ts
JOIN practice p ON ts.practice_id = p.id
LEFT JOIN test_session_answers tsa ON ts.id = tsa.session_id
WHERE ts.user_id = $1
GROUP BY ts.id, p.title, ts.created_at
ORDER BY ts.created_at DESC;

-- name: DeleteById :one
DELETE FROM test_session WHERE id = $1 AND user_id = $2 RETURNING *;