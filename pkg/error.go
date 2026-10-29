package pkg

import "errors"

// ErrNotFound indica que o recurso solicitado não existe.
var ErrNotFound = errors.New("not found")

type ErrorResponse struct {
	Message string `json:"message"`
}
