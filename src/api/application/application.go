package application

import (
	"DX/src/api/middleware"
	"DX/src/utils/logger"
	"fmt"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

var (
	router *gin.Engine
)

func init() {
	router = gin.New()
	zapLogger := logger.GetLogger()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Access-Control-Request-Headers", "Authorization")

	router.Use(cors.New(config))
	router.Use(middleware.RequestLoggerMiddleware())
	router.Use(ginzap.Ginzap(zapLogger, time.RFC3339, true))
	router.MaxMultipartMemory = 8 << 20
}

func StartApplication() {
	initializeRepositories()
	setUpRepositoriesAndManagers()
	mapRoutes()
	//Add custom recovery
	router.Use(gin.CustomRecovery(middleWare.Recovery()))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Error("application start error", router.Run(fmt.Sprintf("%s:%s", "0.0.0.0", port)))
}
