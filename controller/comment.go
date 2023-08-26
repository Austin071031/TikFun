package controller

import (
	"net/http"
	"github.com/RaymondCode/simple-demo/repository"
	"github.com/RaymondCode/simple-demo/service"
	// "github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	// "gorm.io/gorm"
	//"fmt"
)


// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	// token := c.Query("token")
	// decoded_token, err := utils.VerifyToken(token)
	// if err != nil{
	// 	c.JSON(http.StatusOK, repository.Response{StatusCode: 1, StatusMsg: "Verify jwt error"})
	// 	return
	// }
	// users, err := repository.NewUserDaoInstance().QueryUserByName(decoded_token)
	// if len(users) == 0{
	// 	c.JSON(http.StatusOK, repository.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	// 	return
	// }
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
		commentActionResponse,err:=service.CreateComment(content_text,video_id_text,usernameStr)
		if err!=nil{
			//新建评论发生错误，弹出提示窗口
			c.JSON(http.StatusOK,commentActionResponse.Response)
		}else{
      //页面上显示新建的评论
			c.JSON(http.StatusOK,commentActionResponse)
		}
	}else if actionType == "2"{
		//删除
		video_id_text := c.Query("video_id")
		comment_id_text := c.Query("comment_id")
		commentActionResponse,err:=service.DeleteComment(comment_id_text,video_id_text,usernameStr)
		if err!=nil{
      //删除评论发生错误，弹出提示窗口
			c.JSON(http.StatusOK, commentActionResponse.Response)
		}else{
      c.JSON(http.StatusOK, repository.Response{StatusCode: 0})
    }
	}else{
		c.JSON(http.StatusOK, repository.Response{StatusCode: 1, StatusMsg: "actionType not supported"})
	}
}
func CommentList(c *gin.Context) {
  //id, _ := c.Get("userId")
	video_id_text := c.Query("video_id")
	commentListResponse,err:=service.CommentList(video_id_text)
  if err!=nil{
			//显示评论发生错误，弹出提示窗口
			c.JSON(http.StatusOK,commentListResponse.Response)
		}else{
      //页面上显示所有的评论
			c.JSON(http.StatusOK, commentListResponse)
		}
}
