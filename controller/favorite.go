package controller

import (
	"github.com/RaymondCode/simple-demo/service"
  // "github.com/RaymondCode/simple-demo/utils"
	// "github.com/RaymondCode/simple-demo/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)



func FavoriteAction(c *gin.Context) {
	// Token := c.Query("token")
	

 //  //验证用户是否存在
	// username, err := utils.VerifyToken(Token)
	// if err != nil{
	// 	c.JSON(http.StatusOK, repository.Response{StatusCode: 1, StatusMsg: "Verify jwt error"})
	// 	return
	// }
	// users, err := repository.NewUserDaoInstance().QueryUserByName(username)
	// if  len(users) == 0{
	// 	c.JSON(http.StatusOK, repository.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	// 	return
	// }
  username, ok := c.Get("username")
  if !ok{
    username = ""
  }
  usernameStr, _ := username.(string)
  
  VideoID := c.Query("video_id")
	ActionType := c.Query("action_type")
	videoIdint, _ := strconv.Atoi(VideoID)
	actiontypeint, _ := strconv.Atoi(ActionType)

	response, _ := service.FavouriteOrNot(usernameStr, videoIdint, actiontypeint)
	c.JSON(http.StatusOK, response)
}

func FavoriteList(c *gin.Context) {
	username, ok := c.Get("username")
  if !ok{
    username = ""
  }
  usernameStr, _ := username.(string)
  
	// username, err := utils.VerifyToken(Token)
 //  if err != nil{
 //  		c.JSON(http.StatusOK, repository.Response{StatusCode: 1, StatusMsg: "Verify jwt error"})
 //  	}
	favoritelist_response := service.FavoriteList(usernameStr)
  
  c.JSON(http.StatusOK, favoritelist_response)
}