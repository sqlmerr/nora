package error

import "net/http"

type APIError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (err APIError) Error() string {
	return err.Message
}

func New(message string, status int) APIError {
	return APIError{
		Message: message,
		Status:  status,
	}
}

var (
	ErrNotFound           = New("not found", http.StatusNotFound)
	ErrInvalidCredentials = New("invalid credentials", http.StatusUnauthorized)
	ErrUsernameOccupied   = New("username occupied", http.StatusConflict)
	ErrInvalidBody        = New("invalid body", http.StatusUnprocessableEntity)
	ErrUnauthorized       = New("unauthorized", http.StatusUnauthorized)
	ErrInvalidToken       = New("invalid token", http.StatusForbidden)
)
