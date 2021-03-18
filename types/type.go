package types

import (
	"command/models"

	"github.com/gbrlsnchs/jwt/v3"
)

type Payload struct {
	jwt.Payload
	Access bool        `json:"access"`
	User   models.User `json:"user"`
}

type TokenRefreshRequest struct {
	UUID string `json:"uuid"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UUID     string `json:"uuid"`
}
