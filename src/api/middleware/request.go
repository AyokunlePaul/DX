package middleware

import (
	"DX/src/utils/logger"
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
)

func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBuffer, _ := io.ReadAll(c.Request.Body)
		firstReader := io.NopCloser(bytes.NewBuffer(requestBuffer))
		secondReader := io.NopCloser(bytes.NewBuffer(requestBuffer)) //We have to create a new Buffer, because firstReader will be read.

		logger.Info(readBody(firstReader))

		c.Request.Body = secondReader
		c.Next()
	}
}

func readBody(reader io.Reader) string {
	readerBuffer := new(bytes.Buffer)
	_, _ = readerBuffer.ReadFrom(reader)

	content := readerBuffer.String()
	return content
}
