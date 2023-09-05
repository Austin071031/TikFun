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
	
  username, ok := c.Get("username")
  if !ok{
    username = ""
  }
  usernameStr, _ := username.(string)
  
  VideoID := c.Query("video_id")
	ActionType := c.Query("action_type")
	videoIdint, _ := strconv.Atoi(VideoID)
	actiontypeint, _ := strconv.Atoi(ActionType)

	response := service.FavouriteOrNot(usernameStr, videoIdint, actiontypeint)
	c.JSON(http.StatusOK, response)
}

func FavoriteList(c *gin.Context) {
	username, ok := c.Get("username")
  if !ok{
    username = ""
  }
  usernameStr, _ := username.(string)
  
	favoritelist_response := service.FavoriteList(usernameStr)
  
  c.JSON(http.StatusOK, favoritelist_response)
}