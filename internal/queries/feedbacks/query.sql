-- name: CreateAIFeedback :one
UPDATE test_session SET ai_feedback = $1::JSONB WHERE id = $2 AND user_id = $3 RETURNING *;
