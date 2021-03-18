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
