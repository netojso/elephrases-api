package dto

type CreateDeckDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Visibility  string `json:"visibility"`
}

type UpdateDeckDTO = CreateDeckDTO
