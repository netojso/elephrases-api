-- Fila de cartões vencidos: WHERE deck_id = ? AND next_review_at <= ?
CREATE INDEX idx_flashcards_deck_next_review ON flashcards (deck_id, next_review_at);

-- Cartões novos (last_review_at IS NULL) e revisados hoje (last_review_at > ?)
CREATE INDEX idx_flashcards_deck_last_review ON flashcards (deck_id, last_review_at);

ALTER TABLE flashcards
    ADD CONSTRAINT chk_flashcards_state CHECK (state IN ('new', 'learning', 'review', 'lapsed'));

-- Cartões criados antes desta migration foram gravados com ease_factor = 0.
UPDATE flashcards SET ease_factor = 2.5 WHERE ease_factor < 1.3;

ALTER TABLE flashcards
    ADD CONSTRAINT chk_flashcards_ease CHECK (ease_factor >= 1.3);

ALTER TABLE flashcards DROP CONSTRAINT fk_deck;
ALTER TABLE flashcards
    ADD CONSTRAINT fk_deck FOREIGN KEY (deck_id) REFERENCES decks(id) ON DELETE CASCADE;

-- Nullable por enquanto, para não quebrar decks existentes.
ALTER TABLE decks ADD COLUMN user_id UUID REFERENCES users(id);
