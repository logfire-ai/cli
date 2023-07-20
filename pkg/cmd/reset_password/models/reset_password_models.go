package models

type ResetPasswordRequest struct {
	Password string `json:"password"`
}

type ResetPasswordResponse struct {
	IsSuccessful bool   `json:"isSuccessful"`
	Message      string `json:"message,omitempty"`
}
