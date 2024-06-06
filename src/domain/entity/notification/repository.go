package notification

type writer interface {
	SendNotification(Notification) error
}

type reader interface {
	GetAllNotifications(string) ([]Notification, error)
}

type Repository interface {
	reader
	writer
}
