-- name: GetByModuleId :many
SELECT id, number, domain, difficulty, paragraph, prompt, explanation, correct, section_id
FROM question 
WHERE section_id = $1
ORDER BY number;

-- name: GetByModuleWithChoices :many
SELECT
    q.id AS question_id,
    q.number,
    q.domain,
    q.difficulty,
    q.svg,
    q.paragraph,
    q.prompt,
    q.explanation,
    q.correct,
    q.section_id,
    a.id AS answer_id,
    a.label AS answer_label,
    a.text AS answer_text
FROM question q
LEFT JOIN answer_choice a ON a.question_id = q.id
WHERE q.section_id = $1
ORDER BY q.number, a.label;



-- name: CreateWithAnswerChoices :one
WITH new_question AS (
    INSERT INTO question (domain, number, section_id, paragraph, correct, svg, prompt, explanation, difficulty)
    VALUES (
        $1,  -- domain
        $2,  -- number
        $3,  -- section_id
        $4,  -- paragraph
        $5,  -- correct
        $6,  -- svg (can be NULL)
        $7,  -- prompt
        $8,  -- explanation
        $9   -- difficulty
    )
    RETURNING id
)
INSERT INTO answer_choice (question_id, label, text)
VALUES
    ((SELECT id FROM new_question), $10, $11),
    ((SELECT id FROM new_question), $12, $13),
    ((SELECT id FROM new_question), $14, $15),
    ((SELECT id FROM new_question), $16, $17)
RETURNING *;
