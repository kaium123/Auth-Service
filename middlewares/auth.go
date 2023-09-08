package middlewares

import (
	"auth/common/logger"
	"auth/common/utils"
	"auth/config"
	"auth/redis"
	"auth/repository"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string
		cookie, err := ctx.Cookie("access_token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[1]
		} else if err == nil {
			accessToken = cookie
		}

		if accessToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		sub, err := utils.ValidateToken(accessToken, config.Config.AccessTokenPublicKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		raventClient := logger.NewRavenClient()
		logger := logger.NewLogger(raventClient)
		redisConn := redis.NewRedisDb()
		redisRepo := repository.NewRedisRepository(redisConn, logger)
		data, err := redisRepo.Get(context.Background(), "access_token:"+accessToken)
		if err != nil {
			fmt.Println(err)
		}
		if data == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "unauthorized user"})
			return
		}

		ctx.Set("user_id", int64(sub.(float64)))
		ctx.Next()
	}
}
