package models

type UserLogin struct {
	ID       int
	Login    string `json:"login"`
	Password string `json:"password"`
}
