package models

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type RefreshToken struct {
	Token string `json:"token"`
	UUID string `json:"uuid"`
}
