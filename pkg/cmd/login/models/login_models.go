package models

type UserBody struct {
	ProfileID string `json:"profileId"`
	TeamID    string `json:"teamId"`
	Onboarded bool   `json:"onboarded"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

type BearerToken struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Exp          string `json:"exp"`
	Iat          string `json:"iat"`
}

type Response struct {
	IsSuccessful bool        `json:"isSuccessful"`
	Code         int         `json:"code"`
	Email        string      `json:"email"`
	UserBody     UserBody    `json:"userBody"`
	BearerToken  BearerToken `json:"bearerToken"`
	Message      []string    `json:"message"`
}

type SigninRequest struct {
	Email      string `json:"email,omitempty"`
	AuthType   int    `json:"authType"`
	Credential string `json:"credential"`
}
