ALTER TABLE decks DROP COLUMN user_id;

ALTER TABLE flashcards DROP CONSTRAINT fk_deck;
ALTER TABLE flashcards
    ADD CONSTRAINT fk_deck FOREIGN KEY (deck_id) REFERENCES decks(id);

ALTER TABLE flashcards DROP CONSTRAINT chk_flashcards_ease;
ALTER TABLE flashcards DROP CONSTRAINT chk_flashcards_state;

DROP INDEX idx_flashcards_deck_last_review;
DROP INDEX idx_flashcards_deck_next_review;
