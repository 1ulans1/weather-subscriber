package ports

type EmailPublisher interface {
	PublishConfirm(email, token string) error
	PublishNotification(email, city, temperature, condition, unsubToken string) error
}
