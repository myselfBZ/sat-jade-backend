-- name: Create :one
INSERT INTO ai_feedbacks (
    result_id, 
    user_id, 
    content
) VALUES (
    $1, $2, $3
) 
RETURNING id, result_id, user_id, content, created_at;
