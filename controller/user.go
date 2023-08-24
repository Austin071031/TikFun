package controller

import (
  "github.com/RaymondCode/simple-demo/service"
  // "github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	"net/http"
  // "encoding/json"
	// "sync/atomic"
  // "gorm.io/gorm"
  // "log"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
// var usersLoginInfo = map[string]repository.User{
// 	"zhangleidouyin": {
// 		Id:            1,
// 		Name:          "zhanglei",
// 		FollowCount:   10,
// 		FollowerCount: 5,
// 		IsFollow:      true,
// 	},
// }

var userIdSequence = int64(1)

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	userloginresponse, err := service.Register(username, password)
	if err != nil {
		c.JSON(http.StatusOK, userloginresponse)
	}
	c.JSON(http.StatusOK, userloginresponse)
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	userloginresponse, _ := service.Login(username, password)
	c.JSON(http.StatusOK, userloginresponse)
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	userresponse, _ := service.UserInfo(token)
	c.JSON(http.StatusOK, userresponse)
}