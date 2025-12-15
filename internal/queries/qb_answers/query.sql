-- name: Create :one
INSERT INTO question_bank_answers (
    answer,
    user_id,
    question_id,
    status,
    response_duration
) VALUES (
    $1,        
    $2,        
    $3,        
    $4,        
    $5         
) RETURNING *;


-- name: Delete :one
DELETE FROM question_bank_answers WHERE question_id = $1 AND user_id = $2 RETURNING *;


-- name: GetByUser :many
SELECT * FROM question_bank_answers WHERE user_id = $1;
