-- name: GetByQuestionId :many
SELECT *
FROM answer_choice 
WHERE question_id = $1
ORDER BY id;