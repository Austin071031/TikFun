package controller

import (
	"net/http"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)


// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
  username, ok := c.Get("username")
  if !ok{
    username = ""
  }
  usernameStr, _ := username.(string)

  actionType := c.Query("action_type")
	if actionType == "1" {
		//插入
		content_text := c.Query("comment_text")
		video_id_text := c.Query("video_id")
		commentActionResponse := service.CreateComment(content_text,video_id_text,usernameStr)
		c.JSON(http.StatusOK,commentActionResponse)
	}else if actionType == "2"{
		//删除
		video_id_text := c.Query("video_id")
		comment_id_text := c.Query("comment_id")
    commentActionResponse := service.DeleteComment(comment_id_text,video_id_text,usernameStr)
		c.JSON(http.StatusOK, commentActionResponse)
	}
}
func CommentList(c *gin.Context) {
  //id, _ := c.Get("userId")
	video_id_text := c.Query("video_id")
	commentListResponse := service.CommentList(video_id_text)
	c.JSON(http.StatusOK,commentListResponse)
}
