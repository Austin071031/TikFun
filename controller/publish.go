package controller

import (
	"net/http"

	"github.com/RaymondCode/simple-demo/repository"
	"github.com/RaymondCode/simple-demo/service"
	// "github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// 定义了视频列表响应的结构体
//type VideoListResponse struct {
//	Response          //响应参数
//	Token     string  `json:"token"`      //用户令牌
//	VideoList []Video `json:"video_list"` //视频数组
//}

// 视频发布的处理函数 Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")    //从请求参数中获取 token
	title := c.PostForm("title")    //从请求参数中获取 title
	data, err:= c.FormFile("data") //获取上传的文件数据
  if err != nil{
    c.JSON(http.StatusOK, service.VideoListResponse{
      Response: repository.Response{
        StatusCode: 1,
        StatusMsg: "upload video error: " + err.Error(),
      },
    })
  }

	VideoListResponse, _ := service.Publish(token, title, data)
	c.JSON(http.StatusOK, VideoListResponse)
}

// 获取发布视频列表 PublishList all users have same publish video list
func PublishList(c *gin.Context) {

	token := c.Query("token")
  
	VideoListResponse, err := service.PublishList(token)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse)
	}
	c.JSON(http.StatusOK, VideoListResponse)

}
