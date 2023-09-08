package route

import (
	"auth/common/logger"
	"auth/controller"
	"auth/db"
	"auth/middlewares"
	"auth/pb"
	"auth/redis"
	"auth/repository"
	"auth/service"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func Setup() *gin.Engine {
	gin.SetMode(viper.GetString("GIN_MODE"))

	r := gin.New()
	setupCors(r)

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	api := r.Group("/api")

	db := db.InitDB()
	conn, err := grpc.Dial(viper.GetString("ATTACHMENTURL"), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	gRPCCLient := pb.NewAttachmentServiceClient(conn)
	raventClient := logger.NewRavenClient()
	logger := logger.NewLogger(raventClient)
	repo := repository.NewUserRepository(db, logger)
	redisConn := redis.NewRedisDb()
	redisRepo := repository.NewRedisRepository(redisConn, logger)
	service := service.NewUserService(gRPCCLient, repo, redisRepo)
	userController := controller.NewUserController(service)

	auth := api.Group("/auth")

	auth.POST("/login", userController.LogIn)
	auth.POST("/register", userController.Register)

	user := api.Group("/user").Use(middlewares.Auth())
	user.POST("/update/:id", userController.UpdateProfile)
	user.GET("/view/:id", userController.ViewProfile)
	user.POST("/logout",userController.LogOut)

	return r
}

func setupCors(r *gin.Engine) {
	allowConf := viper.GetString("CORS_ALLOW_ORIGINS")
	if allowConf == "" {
		r.Use(cors.Default())
		return
	}
	allowSites := strings.Split(allowConf, ",")
	config := cors.DefaultConfig()
	config.AllowOrigins = allowSites
}
