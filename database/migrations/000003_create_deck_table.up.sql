CREATE TABLE decks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(255) NOT NULL,
    visibility VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL
);
