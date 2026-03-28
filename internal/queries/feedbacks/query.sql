-- name: Create :one
INSERT INTO ai_feedbacks (
    result_id, 
    user_id, 
    header,
    body,
    footer
) VALUES (
    $1, $2, $3, $4, $5
) 
RETURNING *; 

-- name: Get :one
SELECT * FROM ai_feedbacks WHERE result_id = $1 LIMIT 1;
