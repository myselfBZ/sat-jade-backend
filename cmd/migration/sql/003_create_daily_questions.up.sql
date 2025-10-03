CREATE TABLE IF NOT EXISTS daily_question(
    id SERIAL PRIMARY KEY,
    domain VARCHAR(50) NOT NULL,
    paragraph TEXT NOT NULL,
    correct CHAR(1) NOT NULL,
    svg TEXT, -- store file path
    prompt TEXT NOT NULL,
    explanation TEXT NOT NULL,
    difficulty VARCHAR(10) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)