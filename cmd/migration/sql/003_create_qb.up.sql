CREATE TABLE IF NOT EXISTS question_bank (
    id SERIAL PRIMARY KEY,
    domain VARCHAR(255) NOT NULL,
    paragraph TEXT NOT NULL,
    skill   VARCHAR(255),
    question_id VARCHAR(255),
    correct VARCHAR(255) NOT NULL,
    prompt TEXT NOT NULL,
    explanation TEXT NOT NULL,
    difficulty VARCHAR(10) NOT NULL,
    answer_type VARCHAR(255) CHECK (answer_type IN ('mcp', 'oe')), -- Multiple Choice and Open-Ended
    active BOOLEAN DEFAULT FALSE,
    choice_a TEXT,
    choice_b TEXT,
    choice_c TEXT,
    choice_d TEXT
);
