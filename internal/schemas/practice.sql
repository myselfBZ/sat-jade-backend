CREATE TABLE practice (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(10) NOT NULL DEFAULT 'draft' CHECK (status IN ('ready', 'draft'))
);

CREATE TABLE module (
    id SERIAL PRIMARY KEY,
    practice_id INT NOT NULL REFERENCES practice(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE question (
    id SERIAL PRIMARY KEY,
    domain VARCHAR(50) NOT NULL,
    number INT DEFAULT 0,
    section_id INT NOT NULL REFERENCES module(id) ON DELETE CASCADE,
    paragraph TEXT NOT NULL,
    correct CHAR(1) NOT NULL,
    svg TEXT, -- store file path
    prompt TEXT NOT NULL,
    explanation TEXT NOT NULL,
    difficulty VARCHAR(10) NOT NULL
);

CREATE TABLE answer_choice (
    id SERIAL PRIMARY KEY,
    question_id INT NOT NULL REFERENCES question(id) ON DELETE CASCADE,
    label CHAR(1) NOT NULL,
    text TEXT NOT NULL
);


CREATE TABLE test_session (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    practice_id INT NOT NULL REFERENCES practice(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    english_score INT,
    math_score INT,
    total_score INT
);

CREATE TABLE test_session_answers (
    id SERIAL PRIMARY KEY,
    user_answer CHAR(1),
    session_id INT NOT NULL REFERENCES test_session(id) ON DELETE CASCADE,
    correct_answer CHAR(1) NOT NULL,
    module VARCHAR(50) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'omitted' CHECK (status IN ('correct', 'incorrect', 'omitted'))
);