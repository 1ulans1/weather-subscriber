package domain

type ConfirmMessage struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type NotifyMessage struct {
	Email   string `json:"email"`
	City    string `json:"city"`
	Weather struct {
		Temperature string `json:"temperature"`
		Condition   string `json:"condition"`
	} `json:"weather"`
	UnsubToken string `json:"unsub_token"`
}
