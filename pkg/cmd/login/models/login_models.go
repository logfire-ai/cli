package models

type UserBody struct {
	ProfileID string `json:"profileId"`
	AccountID string `json:"accountId"`
	Onboarded bool   `json:"onboarded"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

type TeamBody struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	AccountId string `json:"accountId"`
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
	TeamBody     TeamBody    `json:"teamBody"`
	BearerToken  BearerToken `json:"bearerToken"`
	Message      []string    `json:"message"`
}

type SigninRequest struct {
	Email      string `json:"email,omitempty"`
	AuthType   int    `json:"authType"`
	Credential string `json:"credential"`
}
