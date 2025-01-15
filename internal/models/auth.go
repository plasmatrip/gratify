package models

type LoginRequest struct {
	ID       int32  `json:"-"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}
