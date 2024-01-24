package request

type CreatePlayerRequest struct {
	Name    string `form:"name"`
	Avatar  []byte `form:"avatar"`
	PinCode int    `form:"pin_code"`
}
