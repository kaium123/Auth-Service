package controller

import (
	"auth/common/logger"
	"auth/errors"
	"auth/models"
	"auth/service"
	"net/http"
	"strconv"

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

	if err := c.service.Register(user); err != nil {
		logger.LogError("failed to query user ", err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ginContext.JSON(http.StatusCreated, gin.H{"userId": user.ID, "email": user.Email})

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
