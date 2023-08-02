package models

type SignupRequest struct {
	Email string `json:"email"`
}

type OnboardRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type SetPassword struct {
	Password string `json:"password"`
}

type SignupResponse struct {
	IsSuccessful bool     `json:"isSuccessful"`
	Message      []string `json:"message"`
}
