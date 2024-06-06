package handler

import "github.com/gin-gonic/gin"

type Timeline interface {
	RequestUpdate(*gin.Context)
	PostUpdate(*gin.Context)
	GetTimelineForErrand(*gin.Context)
}

func NewTimelineHandler() Timeline {
	return &timeline{}
}

type timeline struct {
}

func (t *timeline) RequestUpdate(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (t *timeline) PostUpdate(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (t *timeline) GetTimelineForErrand(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}
