package models

type UpdateProfileRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UpdateProfileResponse struct {
	IsSuccessful bool `json:"isSuccessful"`
}
