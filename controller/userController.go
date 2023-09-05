package controller

import (
  "github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

var userIdSequence = int64(1)

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	userloginresponse := service.Register(username, password)
	c.JSON(http.StatusOK, userloginresponse)
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	userloginresponse := service.Login(username, password)
	c.JSON(http.StatusOK, userloginresponse)
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	userresponse := service.UserInfo(token)
	c.JSON(http.StatusOK, userresponse)
}