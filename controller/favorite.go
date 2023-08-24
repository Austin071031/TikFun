package controller

import (
	"github.com/RaymondCode/simple-demo/service"
  "github.com/RaymondCode/simple-demo/repository"
	"github.com/RaymondCode/simple-demo/utils"
	"strconv"
  "net/http"
	"github.com/gin-gonic/gin"
)



func FavoriteAction(c *gin.Context) {
	Token := c.Query("token")
	VideoID := c.Query("video_id")
	ActionType := c.Query("action_type")

  //验证用户是否存在
	username, err := utils.VerifyToken(Token)
	if err != nil{
		c.JSON(http.StatusOK, repository.Response{StatusCode: 1, StatusMsg: "Verify jwt error"})
		return
	}
	_, length, err := repository.NewUserDaoInstance().QueryUserByName(username)
	if  length == 0{
		c.JSON(http.StatusOK, repository.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	videoIdint, _ := strconv.Atoi(VideoID)
	actiontypeint, _ := strconv.Atoi(ActionType)

	response, _ := service.FavouriteOrNot(username, videoIdint, actiontypeint)
	c.JSON(http.StatusOK, response)
}

func FavoriteList(c *gin.Context) {
	Token := c.Query("token")
	

	username, err := utils.VerifyToken(Token)
  if err != nil{
  		c.JSON(http.StatusOK, repository.Response{StatusCode: 1, StatusMsg: "Verify jwt error"})
  	}
	favoritelist_response := service.FavoriteList(username)
  
  c.JSON(http.StatusOK, favoritelist_response)
}