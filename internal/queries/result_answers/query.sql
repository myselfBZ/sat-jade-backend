-- name: CreateMany :many
INSERT INTO result_answers (
    result_id,
    question_id,
    answer_id,
    status
) VALUES (
    unnest(@result_id::int[]),
    unnest(@question_id::int[]),
    NULLIF(unnest(@answer_id::int[]), 0),
    unnest(@status::varchar(10)[])
)
RETURNING *;

-- name: GetByResultID :many
SELECT 
    ra.id,
    ra.result_id,
    ra.question_id,
    ra.status,
    ra.answer_id AS user_answer_id,
    -- Get the label the student picked (e.g., 'B')
    COALESCE(ac_user.label, '') AS user_answer,
    -- Question details
    q.correct AS correct_answer,
    q.number,
    q.paragraph AS passage,
    m.name AS module,
    q.prompt AS question,
    q.explanation,
    -- Pivot the choices into A, B, C, D fields
    MAX(CASE WHEN ac.label = 'A' THEN ac.text END) AS choice_a,
    MAX(CASE WHEN ac.label = 'B' THEN ac.text END) AS choice_b,
    MAX(CASE WHEN ac.label = 'C' THEN ac.text END) AS choice_c,
    MAX(CASE WHEN ac.label = 'D' THEN ac.text END) AS choice_d
FROM result_answers ra
JOIN question q ON ra.question_id = q.id
JOIN module m ON q.section_id = m.id
JOIN answer_choice ac ON ac.question_id = q.id
LEFT JOIN answer_choice ac_user ON ra.answer_id = ac_user.id
WHERE ra.result_id = $1
GROUP BY ra.id, q.id, m.id, ac_user.id
ORDER BY q.number;

-- SELECT 
--     id, 
--     result_id, 
--     question_id, 
--     answer_id, 
--     status 
-- FROM result_answers 
-- WHERE result_id = $1 
-- ORDER BY id;
