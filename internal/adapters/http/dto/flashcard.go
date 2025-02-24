package dto

import "encoding/json"

type CreateFlashcardDTO struct {
	DeckID string `json:"deck_id" binding:"required"`
	Front  string `json:"front" binding:"required"`
	Back   string `json:"back" binding:"required"`
	Media  string `json:"media"`
}

type ReviewFlashcardDTO struct {
	FlashCardID string `json:"flashcard_id" binding:"required"`
	Response    string `json:"response" binding:"required"`
}

func (c *CreateFlashcardDTO) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (r *ReviewFlashcardDTO) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalCreateFlashcardDTO(data []byte) (*CreateFlashcardDTO, error) {
	var c CreateFlashcardDTO
	err := json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
