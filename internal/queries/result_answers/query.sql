-- name: CreateMany :many
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

-- name: GetByResultID :many
SELECT * FROM test_session_answers WHERE session_id = $1 ORDER BY id;