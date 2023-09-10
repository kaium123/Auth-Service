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
		logger.LogError(err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := c.service.Register(user)

	if err != nil {
		logger.LogError(err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ginContext.JSON(http.StatusCreated, gin.H{"userId": id, "email": user.Email})

}

func (c *UserController) LogIn(ginContext *gin.Context) {
	signInInfo := &models.SignInData{}
	if err := ginContext.Bind(&signInInfo); err != nil {
		logger.LogError(err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logger.LogInfo(signInInfo)

	resp, err := c.service.LogIn(*signInInfo)

	if err != nil {
		logger.LogError(err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ginContext.JSON(http.StatusCreated, gin.H{"user": resp})

}

func (c *UserController) UpdateProfile(ginContext *gin.Context) {
	var user models.User
	if err := ginContext.Bind(&user); err != nil {
		logger.LogError(err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := int(ginContext.GetInt64("user_id"))
	user.ID = userID
	err := c.service.UpdateProfile(&user)
	if err != nil {
		logger.LogError(err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ginContext.JSON(http.StatusCreated, gin.H{"user": user})

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

func (c *UserController) MyProfile(ginContext *gin.Context) {

	userID := int(ginContext.GetInt64("user_id"))
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
		logger.LogError(err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = c.service.LogOut(accessToken)
	if err != nil {
		logger.LogError(err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ginContext.JSON(http.StatusCreated, gin.H{"msg": "successfully log out"})

}

func (r *UserController) RequestSent(ginContext *gin.Context) {
	userID := int(ginContext.GetInt64("user_id"))
	requestedIDString := ginContext.Params.ByName("id")
	requestedID, err := strconv.Atoi(requestedIDString)
	if err != nil {
		logger.LogError(err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = r.service.RequestSent((userID), requestedID)

	if err != nil {
		logger.LogError(err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ginContext.JSON(http.StatusCreated, gin.H{"msg": "request sent"})
}

func (r *UserController) RequestAccept(ginContext *gin.Context) {
	userID := int(ginContext.GetInt64("user_id"))

	requestedIDString := ginContext.Params.ByName("id")
	requestedID, err := strconv.Atoi(requestedIDString)
	if err != nil {
		logger.LogError(err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = r.service.RequestAccept(userID, requestedID)

	if err != nil {
		logger.LogError(err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ginContext.JSON(http.StatusCreated, gin.H{"msg": "request accepted"})
}

func (r *UserController) ManageConnection(ginContext *gin.Context) {
	userID := int(ginContext.GetInt64("user_id"))

	friendIDString := ginContext.Params.ByName("id")
	friendID, err := strconv.Atoi(friendIDString)
	if err != nil {
		logger.LogError(err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = r.service.ManageConnection(userID, friendID)
	if err != nil {
		logger.LogError(err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ginContext.JSON(http.StatusCreated, gin.H{"msg": "unfriend successfully"})
}

func (r *UserController) ViewFriends(ginContext *gin.Context) {
	userID := int(ginContext.GetInt64("user_id"))
	logger.LogInfo(userID)

	resp, err := r.service.ViewFriends(userID)
	if err != nil {
		logger.LogError(err)
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ginContext.JSON(http.StatusCreated, resp)
}
