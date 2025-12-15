CREATE TABLE IF NOT EXISTS question_bank_answers(
    answer VARCHAR(255) NOT NULL,
    user_id     UUID NOT NULL REFERENCES users(id),
    question_id INT NOT NULL REFERENCES question_bank(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status VARCHAR(10) NOT NULL DEFAULT 'incorrect' CHECK (status IN ('correct', 'incorrect')),
    response_duration INT
)
