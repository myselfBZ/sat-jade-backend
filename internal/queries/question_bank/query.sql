-- name: Create :one
INSERT INTO question_bank (
    domain,
    skill,
    paragraph,
    question_id,
    correct,
    prompt,
    explanation,
    difficulty,
    choice_a,
    choice_b,
    choice_c,
    choice_d,
    answer_type,
    active

) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
)
RETURNING *;


-- name: GetById :one
SELECT question_bank.*, question_bank_answers.* 
FROM question_bank 
LEFT JOIN question_bank_answers 
ON question_bank.id = question_bank_answers.question_id AND question_bank_answers.user_id = $1
WHERE question_bank.id = $2;

-- name: GetIdBySkill :many
SELECT id FROM question_bank WHERE skill = $1;

-- name: GetCollectionDetails :many
SELECT COUNT(id), domain, skill FROM question_bank GROUP BY domain, skill;
