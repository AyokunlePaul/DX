package application

import (
	"DX/src/api/handler"
	"DX/src/api/handler/admin"
	"DX/src/api/middleware"
	"DX/src/domain/entity/auth"
	"DX/src/domain/entity/category"
	errandRepository "DX/src/domain/entity/errand"
	fileRepository "DX/src/domain/entity/file"
	"DX/src/domain/entity/notification"
	secRepository "DX/src/domain/entity/security"
	"DX/src/domain/entity/user"
	"DX/src/domain/entity/wallet"
	adminUseCase "DX/src/domain/usecase/admin"
	"DX/src/domain/usecase/authentication"
	"DX/src/domain/usecase/errand"
	"DX/src/domain/usecase/file"
	"DX/src/domain/usecase/init_data"
	"DX/src/domain/usecase/security"
	wallet2 "DX/src/domain/usecase/wallet"
	"DX/src/pkg/error_service"
	"DX/src/pkg/password_service"
	"DX/src/pkg/token_service"
	"DX/src/utils/logger"
	"cloud.google.com/go/storage"
	"context"
	ginzap "github.com/gin-contrib/zap"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/option"
	"time"
)

func initializeRepositories() {
}

const (
	mongoUri = "mongodb://localhost:27017"
)

var (
	authenticationHandler handler.Authentication
	securityHandler       handler.Security
	errandHandler         handler.Errand
	initHandler           handler.Init
	walletHandler         handler.Wallet
	categoryHandler       admin.Category
	userAdminHandler      admin.User
	errandAdminHandler    admin.Errand
	middleWare            middleware.Middleware
)

func GetDatabase() *mongo.Database {
	mongoContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mongoClient, clientError := mongo.Connect(mongoContext, options.Client().ApplyURI(mongoUri))
	if clientError != nil {
		panic(clientError)
	}

	return mongoClient.Database("errand-app")
}

func InitializeUserCollection(database *mongo.Database) *mongo.Collection {
	mongoContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{"phone_number", 1},
		},
		Options: options.Index().SetUnique(true),
	}
	emailIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{"email", "text"},
			{"first_name", "text"},
			{"last_name", "text"},
		},
		Options: options.Index(),
	}
	collection := database.Collection("users")
	_, indexError := collection.Indexes().CreateMany(mongoContext, []mongo.IndexModel{userIndexModel, emailIndexModel})
	if indexError != nil {
		panic(indexError)
	}

	return collection
}

func InitializeTokenCollection(database *mongo.Database) *mongo.Collection {
	mongoContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tokenExpIndex := mongo.IndexModel{
		Keys: bson.D{
			{"access_token_expires", 1},
		},
		Options: options.Index().SetExpireAfterSeconds(0),
	}
	uuidIndex := mongo.IndexModel{
		Keys: bson.D{
			{"access_uuid", 1},
		},
	}

	collection := database.Collection("token")
	_, indexError := collection.Indexes().CreateMany(mongoContext, []mongo.IndexModel{tokenExpIndex, uuidIndex})
	if indexError != nil {
		panic(indexError)
	}

	return collection
}

func InitializeSecurityCollection(database *mongo.Database) *mongo.Collection {
	mongoContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userIdIndex := mongo.IndexModel{
		Keys: bson.D{
			{"user_id", 1},
			{"phone_number", 1},
		},
		Options: options.Index().SetUnique(true),
	}

	collection := database.Collection("security-questions")
	_, indexError := collection.Indexes().CreateMany(mongoContext, []mongo.IndexModel{userIdIndex})
	if indexError != nil {
		panic(indexError)
	}

	return collection
}

func InitializeErrandCollection(database *mongo.Database) *mongo.Collection {
	mongoContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indices := mongo.IndexModel{
		Keys: bson.D{
			{"_id", 1},
			{"user_id", 1},
		},
		Options: options.Index(),
	}

	collection := database.Collection("errands")
	_, indexError := collection.Indexes().CreateMany(mongoContext, []mongo.IndexModel{indices})
	if indexError != nil {
		panic(indexError)
	}

	return collection
}

func InitializeCategoryCollection(database *mongo.Database) *mongo.Collection {
	mongoContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indices := mongo.IndexModel{
		Keys: bson.D{
			{"_id", 1},
		},
		Options: options.Index(),
	}
	uniqueIndices := mongo.IndexModel{
		Keys: bson.D{
			{"_id", 1},
			{"identifier", 1},
		},
		Options: options.Index().SetUnique(true),
	}

	collection := database.Collection("categories")
	_, indexError := collection.Indexes().CreateMany(mongoContext, []mongo.IndexModel{indices, uniqueIndices})
	if indexError != nil {
		panic(indexError)
	}

	return collection
}

func InitializeNotificationCollection(database *mongo.Database) *mongo.Collection {
	mongoContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indices := mongo.IndexModel{
		Keys: bson.D{
			{"_id", 1},
		},
		Options: options.Index(),
	}

	collection := database.Collection("notifications")
	_, indexError := collection.Indexes().CreateMany(mongoContext, []mongo.IndexModel{indices})
	if indexError != nil {
		panic(indexError)
	}

	return collection
}

func InitializeTransactionCollection(database *mongo.Database) *mongo.Collection {
	mongoContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indices := mongo.IndexModel{
		Keys: bson.D{
			{"_id", 1},
			{"user_id", 1},
		},
		Options: options.Index(),
	}

	collection := database.Collection("transactions")
	_, indexError := collection.Indexes().CreateMany(mongoContext, []mongo.IndexModel{indices})
	if indexError != nil {
		panic(indexError)
	}

	return collection
}

func setUpRepositoriesAndManagers() {
	//Service
	tokenService := token_service.New()
	errorService := error_service.New()
	passwordService := password_service.New()

	//DB and collection
	db := GetDatabase()
	userCollection := InitializeUserCollection(db)
	tokenCollection := InitializeTokenCollection(db)
	secCollection := InitializeSecurityCollection(db)
	errandCollection := InitializeErrandCollection(db)
	categoryCollection := InitializeCategoryCollection(db)
	notificationCollection := InitializeNotificationCollection(db)
	transactionCollection := InitializeTransactionCollection(db)

	//Clients
	strClient, err := storage.NewClient(context.Background(), option.WithCredentialsFile(""))
	if err != nil {
		logger.Error("Storage Bucket::", err)
	}

	// Repositories
	authRepo := auth.NewMongoCacheRepository(tokenCollection)
	secRepo := secRepository.NewDatabaseRepository(secCollection)
	userRepo := user.NewDatabaseRepository(userCollection)
	errandRepo := errandRepository.NewRepository(errandCollection)
	fileRepo := fileRepository.NewRepository(strClient)
	categoryRepo := category.NewRepository(categoryCollection)
	notificationRepo := notification.NewInAppNotificationRepository(notificationCollection)
	walletRepo := wallet.NewWalletRepository(transactionCollection)

	// Managers
	authManager := auth.NewManager(tokenService, authRepo)

	// UseCases
	authUseCase := authentication.NewUseCase(userRepo, errorService, passwordService, authManager, notificationRepo)
	secUseCase := security.NewUseCase(authManager, secRepo, errorService, userRepo)
	errandUseCase := errand.NewUseCase(authManager, errandRepo, userRepo, errorService, notificationRepo, categoryRepo, errandRepo, walletRepo)
	fileUseCase := file.NewUseCase(authManager, fileRepo, errandRepo, errorService)
	adminUserUseCase := adminUseCase.NewUserUseCase(authManager, userRepo, errorService)
	adminCategoryUseCase := adminUseCase.NewCategoryUseCase(fileRepo, authManager, userRepo, categoryRepo, errorService)
	adminErrandUseCase := adminUseCase.NewErrandUseCase(authManager, errandRepo, errorService, userRepo, notificationRepo, categoryRepo)
	initUseCase := init_data.NewUseCase(categoryRepo)
	walletUseCase := wallet2.NewUseCase(walletRepo, errorService, authManager)

	// Middlewares
	middleWare = middleware.NewErrandMiddleware(userRepo, tokenService, authManager)

	// Handlers
	authenticationHandler = handler.NewAuthHandler(authUseCase, errandUseCase)
	securityHandler = handler.NewSecurityHandler(secUseCase)
	errandHandler = handler.NewErrandHandler(errandUseCase, fileUseCase)
	initHandler = handler.NewInitHandler(initUseCase)
	userAdminHandler = admin.NewAdminUserHandler(adminUserUseCase)
	categoryHandler = admin.NewAdminCategoryHandler(adminCategoryUseCase)
	errandAdminHandler = admin.NewAdminErrandHandler(adminErrandUseCase, errandUseCase)
	walletHandler = handler.NewWalletHandler(walletUseCase)

	zapLogger := logger.GetLogger()

	router.Use(middleWare.CORS())
	router.Use(middleware.RequestLoggerMiddleware())
	router.Use(ginzap.Ginzap(zapLogger, time.RFC3339, true))
}

func mapRoutes() {
	v1Group := router.Group("/v1")
	{
		v1Group.GET("/init", middleWare.Authorization(), middleWare.Suspension(), initHandler.GetInitData)

		v1Group.POST("/security-question", middleWare.Authorization(), middleWare.Suspension(), securityHandler.SaveSecurityQuestion)
		v1Group.GET("/security-question", securityHandler.GetSecurityQuestion)
		v1Group.POST("/security-question/verify", securityHandler.VerifySecurityQuestion)
		v1Group.POST("/paystack/webhook", walletHandler.PaystackWebhook)
		v1Group.POST("/transact", walletHandler.MakePayment)
		v1Group.GET("/errand/market", errandHandler.FetchAllErrands)

		authenticationGroup := v1Group.Group("/user")
		{
			authenticationGroup.POST("/sign-in", authenticationHandler.Login)
			authenticationGroup.POST("/verify-phone", authenticationHandler.VerifyPhoneNumber)
			authenticationGroup.POST("/sign-up", authenticationHandler.CreateAccount)
			authenticationGroup.POST("/sign-out", middleWare.Authorization(), authenticationHandler.LogOut)
			authenticationGroup.POST("/suspend", middleWare.Authorization(), authenticationHandler.SuspendUser)
			authenticationGroup.POST("/password", authenticationHandler.UpdatePassword)
			authenticationGroup.GET("/profile", middleWare.Authorization(), middleWare.Suspension(), authenticationHandler.Profile)
			authenticationGroup.GET("/errands", middleWare.Authorization(), middleWare.Suspension(), authenticationHandler.MyErrands)
			authenticationGroup.GET("/notifications", middleWare.Authorization(), middleWare.Suspension(), authenticationHandler.MyNotifications)
			authenticationGroup.GET("/wallet", middleWare.Authorization(), middleWare.Suspension(), walletHandler.GetWallet)
			authenticationGroup.GET("/:id", middleWare.Authorization(), middleWare.Suspension(), authenticationHandler.GetUser)
			authenticationGroup.POST("/rate", middleWare.Authorization(), middleWare.Suspension(), authenticationHandler.RateUser)
		}
		errandGroup := v1Group.Group("/errand", middleWare.Authorization(), middleWare.Suspension())
		{
			errandGroup.POST("/files", errandHandler.UploadErrandFiles)
			errandGroup.GET("/draft", errandHandler.GetDraftErrand)
			errandGroup.PATCH("/:id", errandHandler.UpdateErrand)
			errandGroup.POST("/:id", errandHandler.CreateErrand)
			errandGroup.DELETE("/:id/cancel", errandHandler.CancelErrand)
			errandGroup.PATCH("/:id/complete", errandHandler.CompleteErrand)
			errandGroup.GET("/:id", errandHandler.GetErrand)
			errandGroup.GET("/categories", categoryHandler.GetAllCategories)
			errandGroup.POST("/start", errandHandler.StartErrand)

			bidGroup := errandGroup.Group("/:id")
			{
				bidGroup.POST("/bid", errandHandler.BidForErrand)
				bidGroup.PUT("/bid/:bid_id", errandHandler.UpdateBidForErrand)
				bidGroup.PUT("/bid/:bid_id/respond", errandHandler.RespondToBid)
				bidGroup.DELETE("/bid/:bid_id/respond", errandHandler.RejectErrandContract)
			}
			timelineGroup := errandGroup.Group("/timeline")
			{
				timelineGroup.POST("", errandHandler.PostUpdate)
				timelineGroup.GET("/request/:id", errandHandler.RequestForUpdate)
			}
		}
		adminGroup := v1Group.Group("/admin", middleWare.Authorization(), middleWare.Suspension(), middleWare.Admin())
		{
			adminGroup.GET("/users", userAdminHandler.GetAllUsers)
			adminGroup.GET("/errands", errandAdminHandler.GetAllErrands)
			userGroup := adminGroup.Group("/user")
			{
				userGroup.GET("/:id", userAdminHandler.GetUser)
				userGroup.POST("", userAdminHandler.CreateUser)
				userGroup.PATCH("", userAdminHandler.UpdateUser)
				userGroup.PATCH("/restore/:id", userAdminHandler.RestoreUser)
				userGroup.PUT("/suspend/:id", userAdminHandler.SuspendUser)
				userGroup.DELETE("/:id", userAdminHandler.DeleteUser)
				userGroup.DELETE("", userAdminHandler.DeleteUsers)
				userGroup.PUT("", userAdminHandler.SuspendUsers)
			}
			categoryGroup := adminGroup.Group("/category")
			{
				categoryGroup.POST("", categoryHandler.CreateCategory)
				categoryGroup.PUT("/:id", categoryHandler.UpdateCategory)
				categoryGroup.GET("", categoryHandler.GetAllCategories)
				categoryGroup.DELETE("/:id", categoryHandler.DeleteCategory)
			}
			errandAdminGroup := adminGroup.Group("/errand")
			{
				errandAdminGroup.POST("", errandAdminHandler.CreatErrand)
				errandAdminGroup.GET("/:id", errandAdminHandler.GetErrand)
			}
		}
	}
}
