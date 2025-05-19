package http

type SubscribeRequest struct {
	Email     string `json:"email" form:"email"`
	City      string `json:"city" form:"city"`
	Frequency string `json:"frequency" form:"frequency"`
}
