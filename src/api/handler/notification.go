package handler

import "github.com/gin-gonic/gin"

type Notification interface {
	FetchAllNotifications(*gin.Context)
}

type notificationImpl struct {
}

func NewNotificationHandler() {

}
