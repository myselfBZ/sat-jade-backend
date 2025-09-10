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



-- name: GetPracticePreviews :many
SELECT * FROM practice;

-- name: GetById :one
SELECT id, title, created_at
FROM practice 
WHERE id = $1;

-- name: GetModulesByPracticeId :many
SELECT id, name, practice_id
FROM module 
WHERE practice_id = $1
ORDER BY id;

-- name: GetQuestionsByModuleId :many
SELECT id, number, domain, difficulty, paragraph, prompt, explanation, correct, section_id
FROM question 
WHERE section_id = $1
ORDER BY number;

-- name: GetAnswerChoicesByQuestionId :many
SELECT id, label, text, question_id
FROM answer_choice 
WHERE question_id = $1
ORDER BY id;

-- name: Delete :one
DELETE FROM practice WHERE id = $1 RETURNING *;

-- name: AddQuestion :one
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

-- name: GetModuleID :one 
SELECT id from module WHERE practice_id = $1 AND name = $2;

-- name: CreateTestSession :one
INSERT INTO test_session(
    practice_id, 
    user_id, 
    english_score, 
    math_score, 
    total_score
    ) VALUES($1, $2, $3, $4, $5) RETURNING *;

-- name: CreateTestSessionAnswers :many
INSERT INTO test_session_answers (
    user_answer,
    session_id,
    correct_answer,
    module,
    status
) VALUES (
    unnest(@user_answer::char(1)[]),
    unnest(@session_id::int[]),
    unnest(@correct_answer::char(1)[]),
    unnest(@module::varchar(50)[]),
    unnest(@status::varchar(10)[])
)
RETURNING *;

-- name: GetExamResultsByUserID :many
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


-- name: GetSessionAnswers :many
SELECT * FROM test_session_answers WHERE session_id = $1 ORDER BY id;

-- name: GetSessionById :one
SELECT * FROM test_session WHERE id = $1;


-- name: DeleteSessionById :one
DELETE FROM test_session WHERE id = $1 AND user_id = $2 RETURNING *;

-- name: GetLastSession :one
SELECT * 
FROM test_session 
WHERE user_id = $1
ORDER BY created_at DESC 
LIMIT 1;

-- name: CreateAIFeedback :one
UPDATE test_session SET ai_feedback = $1::JSONB WHERE id = $2 AND user_id = $3 RETURNING *;

-- name: GetAllSessions :many
SELECT * from test_session;