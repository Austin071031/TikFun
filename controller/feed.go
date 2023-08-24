package controller

import (
	"github.com/RaymondCode/simple-demo/service"
  "github.com/RaymondCode/simple-demo/repository"
	"github.com/gin-gonic/gin"
  "github.com/RaymondCode/simple-demo/utils"
	"net/http"
)

// Feed same demo video list for every request
func Feed(c *gin.Context) {
  token := c.Query("token")
  if token == ""{
    FeedResponse, _ := service.Feed("")
	  c.JSON(http.StatusOK, FeedResponse)
  }
  
	username, err := utils.VerifyToken(token)
	if err != nil{
		c.JSON(http.StatusOK, repository.Response{StatusCode: 1, StatusMsg: "Verify jwt error"})
		return
	}
  FeedResponse, _ := service.Feed(username)
	c.JSON(http.StatusOK, FeedResponse)
	
}
