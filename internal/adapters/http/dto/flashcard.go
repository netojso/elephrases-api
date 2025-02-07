package dto

type CreateFlashcardDTO struct {
	DeckID string `json:"deck_id" binding:"required"`
	Front  string `json:"front" binding:"required"`
	Back   string `json:"back" binding:"required"`
}

type ReviewFlashcardDTO struct {
	FlashCardID string `json:"flashcard_id" binding:"required"`
	Response    string `json:"response" binding:"required"`
}
