package middleware

import (
	"DX/src/domain/entity/auth"
	"DX/src/domain/entity/user"
	"DX/src/pkg/response"
	"DX/src/pkg/token_service"
	"DX/src/utils/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Middleware interface {
	Suspension() gin.HandlerFunc
	Authorization() gin.HandlerFunc
	Admin() gin.HandlerFunc
	SuperAdmin() gin.HandlerFunc
	CORS() gin.HandlerFunc
	Recovery() gin.RecoveryFunc
}

type middleWare struct {
	user.Repository
	token_service.Service
	auth.Manager
}

func NewErrandMiddleware(repo user.Repository, service token_service.Service, manager auth.Manager) Middleware {
	return &middleWare{
		Repository: repo,
		Service:    service,
		Manager:    manager,
	}
}

func (m *middleWare) Suspension() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if strings.TrimSpace(authHeader) == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.NewUnAuthorizedError())
			return
		}
		token := strings.Split(authHeader, " ")[1]
		userId, resp := m.Manager.Get(token)
		if resp != nil {
			ctx.JSON(resp.Code, resp)
			ctx.Abort()
			return
		}

		nUser, err := m.Repository.GetWithId(*userId)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, response.NewUnAuthorizedError())
			ctx.Abort()
			return
		}

		if nUser.IsSuspended || nUser.IsDeleted {
			ctx.JSON(http.StatusUnauthorized, response.NewUnAuthorizedError())
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (m *middleWare) Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenPayload := ctx.GetHeader("Authorization")
		if strings.TrimSpace(tokenPayload) == "" {
			ctx.JSON(http.StatusUnauthorized, response.NewUnAuthorizedError())
			ctx.Abort()
			return
		}
		bearerToken := strings.Split(tokenPayload, " ")
		if len(bearerToken) != 2 {
			ctx.JSON(http.StatusUnauthorized, response.NewUnAuthorizedError())
			ctx.Abort()
			return
		}
		userToken := bearerToken[1]
		if strings.TrimSpace(userToken) == "" {
			ctx.JSON(http.StatusUnauthorized, response.NewUnAuthorizedError())
			ctx.Abort()
			return
		}
		if err := m.Service.CheckTokenValidity(userToken, true); err != nil {
			ctx.JSON(http.StatusUnauthorized, response.NewUnAuthorizedError())
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func (m *middleWare) Admin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if strings.TrimSpace(authHeader) == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.NewUnAuthorizedError())
			return
		}
		token := strings.Split(authHeader, " ")[1]
		userId, resp := m.Manager.Get(token)
		if resp != nil {
			ctx.AbortWithStatusJSON(resp.Code, resp)
			return
		}

		nUser, err := m.Repository.GetWithId(*userId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
			return
		}

		if !nUser.IsAdmin() {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.NewUnAuthorizedError())
			return
		}
		ctx.Next()
	}
}

func (m *middleWare) SuperAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if strings.TrimSpace(authHeader) == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.NewUnAuthorizedError())
			return
		}
		token := strings.Split(authHeader, " ")[1]
		userId, resp := m.Manager.Get(token)
		if resp != nil {
			ctx.AbortWithStatusJSON(resp.Code, resp)
			return
		}

		nUser, err := m.Repository.GetWithId(*userId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
			return
		}

		if !nUser.IsAdmin() {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.NewUnAuthorizedError())
			return
		}
		ctx.Next()
	}
}

func (m *middleWare) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func (m *middleWare) Recovery() gin.RecoveryFunc {
	return func(c *gin.Context, err any) {
		logger.Error("ServerError", err.(error))
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.NewInternalServerError("It's not you, it's us. We are working on a fix.⚒"))
	}
}
