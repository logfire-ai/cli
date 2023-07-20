package logout

type SignOutResponse struct {
	IsSuccessful bool   `json:"isSuccessful"`
	Message      string `json:"message,omitempty"`
}
