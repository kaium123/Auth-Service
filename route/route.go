
package route

import (
	"auth/common/logger"
	"auth/controller"
	"auth/db"
	"auth/repository"
	"auth/service"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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

	raventClient := logger.NewRavenClient()
	logger := logger.NewLogger(raventClient)
	repo := repository.NewUserRepository(db, logger)
	service := service.NewUserService(repo)
	userController := controller.NewUserController(service)

	user := api.Group("/user")

	user.POST("/login", userController.LogIn)
	user.POST("/register", userController.Register)
	user.POST("/update/:id", userController.UpdateProfile)
	user.GET("/view/:id", userController.ViewProfile)


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
