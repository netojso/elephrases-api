CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE flashcards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    front TEXT NOT NULL,
    back TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_review_at TIMESTAMP
    );
