package models

type UpdateProfileRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Role      string `json:"role"`
}

type UpdateProfileResponse struct {
	IsSuccessful bool `json:"isSuccessful"`
}
