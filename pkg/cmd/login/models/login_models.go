package models

type UserBody struct {
	ProfileID string `json:"profileId"`
	TeamID    string `json:"teamId"`
	Onboarded bool   `json:"onboarded"`
	Email     string `json:"email"`
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
