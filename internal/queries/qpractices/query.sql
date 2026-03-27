-- name: Create :one
WITH new_practice AS (
    INSERT INTO practice (title)
    VALUES ($1)
    RETURNING id, title, created_at, status
),
inserted_modules AS (
    INSERT INTO module (practice_id, name)
    SELECT id, name
    FROM new_practice, (VALUES 
        ('Reading And Writing 1'),
        ('Reading And Writing 2'),
        ('Math 1'),
        ('Math 2')
    ) AS t(name)
    RETURNING id
)
SELECT * FROM new_practice;


-- name: GetCorrectAnswersWithChoices :many
SELECT 
    q.id AS question_id,
    q.correct AS correct_label,
    -- Get all 4 choices for this question as a JSON array
    (
        SELECT json_agg(json_build_object(
            'id', ac.id, 
            'label', ac.label
        ) ORDER BY ac.label)
        FROM answer_choice ac 
        WHERE ac.question_id = q.id
    ) AS choices
FROM module m
JOIN question q ON q.section_id = m.id
WHERE m.practice_id = $1
ORDER BY m.id, q.number;

-- name: GetPracticePreviews :many
SELECT * FROM practice;


-- name: Delete :one
DELETE FROM practice WHERE id = $1 RETURNING *;

-- name: GetFullPracticeTest :one
SELECT 
    p.id, 
    p.title, 
    p.created_at,
    (
        SELECT json_agg(m_data)
        FROM (
            SELECT 
                m.id, 
                m.name,
                (
                    SELECT json_agg(q_data)
                    FROM (
                        SELECT 
                            q.id, q.number, q.domain, q.difficulty, q.svg, 
                            q.paragraph, q.prompt, q.explanation, q.correct,
                            json_agg(json_build_object(
                                'id', a.id, 
                                'label', a.label, 
                                'text', a.text
                            ) ORDER BY a.label) AS choices
                        FROM question q
                        LEFT JOIN answer_choice a ON a.question_id = q.id
                        WHERE q.section_id = m.id
                        GROUP BY q.id
                        ORDER BY q.number
                    ) q_data
                ) AS questions
            FROM module m
            WHERE m.practice_id = p.id
            ORDER BY m.id
        ) m_data
    ) AS sections
FROM practice p
WHERE p.id = $1;
