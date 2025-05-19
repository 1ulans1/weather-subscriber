package ports

type EmailSender interface {
	SendConfirmation(to, token string) error
	SendNotification(to, city, temperature, condition, unsubscribeURL string) error
}
