package controller

import (
	"auth/common/logger"
	"auth/errors"
	"auth/models"
	"auth/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	errors.GinError
	service service.UserServiceInterface
}

func NewUserController(service service.UserServiceInterface) *UserController {
	return &UserController{service: service}
}

func (c *UserController) Register(ginContext *gin.Context) {
	var user models.User
	if err := ginContext.Bind(&user); err != nil {
		logger.LogError("failed to query estimate ", err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := c.service.Register(user)

	if err != nil {
		logger.LogError("failed to query user ", err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ginContext.JSON(http.StatusCreated, gin.H{"userId": id, "email": user.Email})

}

func (c *UserController) LogIn(ginContext *gin.Context) {
	signInInfo := &models.SignInData{}
	if err := ginContext.Bind(&signInInfo); err != nil {
		logger.LogError("failed to query estimate ", err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logger.LogInfo(signInInfo)

	resp, err := c.service.LogIn(*signInInfo)

	if err != nil {
		logger.LogError("failed to query estimate ", err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ginContext.JSON(http.StatusCreated, gin.H{"user": resp})

}

func (c *UserController) UpdateProfile(ginContext *gin.Context) {
	var user models.User
	if err := ginContext.Bind(&user); err != nil {
		logger.LogError("failed to query estimate ", err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDString := ginContext.Params.ByName("id")
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		logger.LogError(err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	user.ID = userID
	err = c.service.UpdateProfile(&user)
	if err != nil {
		logger.LogError("failed to query estimate ", err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//ginContext.JSON(http.StatusCreated, gin.H{"user": resp})

}

func (c *UserController) ViewProfile(ginContext *gin.Context) {

	userIDString := ginContext.Params.ByName("id")
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		logger.LogError(err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	profile, err := c.service.ViewProfile(userID)
	if err != nil {
		logger.LogError("failed to query estimate ", err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ginContext.JSON(http.StatusCreated, gin.H{"user": profile})

}

func (c *UserController) LogOut(ginContext *gin.Context) {

	var accessToken string
	cookie, err := ginContext.Cookie("access_token")

	authorizationHeader := ginContext.Request.Header.Get("Authorization")
	fields := strings.Fields(authorizationHeader)

	if len(fields) != 0 && fields[0] == "Bearer" {
		accessToken = fields[1]
	} else if err == nil {
		accessToken = cookie
	}

	if accessToken == "" {
		logger.LogError("failed to query estimate ", err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = c.service.LogOut(accessToken)
	if err != nil {
		logger.LogError("failed to query estimate ", err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ginContext.JSON(http.StatusCreated, gin.H{"msg": "successfully log out"})

}