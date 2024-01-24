package dto

import "time"

type PlayerDTO struct {
	Name         string    `json:"name"`
	Avatar       []byte    `json:"avatar"`
	RegisteredAt time.Time `json:"registered_at"`
	PinCode      int       `json:"pin_code"`
}
