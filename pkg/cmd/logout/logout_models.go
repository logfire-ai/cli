package logout

type SignOutResponse struct {
	IsSuccessful bool   `json:"isSuccessful"`
	Msg          string `json:"msg"`
}
