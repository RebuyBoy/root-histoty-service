package model

import (
	"time"
)

type Player struct {
	ID           string
	Name         string
	Avatar       []byte    `db:"avatar"`
	RegisteredAt time.Time `db:"registered_at"`
	PinCode      int       `db:"pin_code"`
}
