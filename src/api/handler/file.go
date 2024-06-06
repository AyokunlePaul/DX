package handler

import "github.com/gin-gonic/gin"

type File interface {
	UploadFile(*gin.Context)
	UploadFiles(*gin.Context)
}
