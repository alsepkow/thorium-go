package thordb

import "github.com/dgrijalva/jwt-go"

type SessionStatus int

const (
	Disconnected = iota
	Connected
	OK
	Interrupted
)

type AccountSession struct {
	Status SessionStatus `json:"sessionStatus"`
	UserID int           `json:"uid"`
	Token  *jwt.Token    `json:"jwt"`
}
