package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/repository"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	// "gorm.io/gorm"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	// "sync/atomic"
	"time"
)

type VideoListResponse struct {
	repository.Response
	VideoList []Video `json:"video_list"`
	Token     string  `json:"token"`
}

var videoIdSequence = int64(0)
var serverDomain string

func GetServerDomain(c *gin.Context) {
	// 获取服务器的域名
	domain := c.Request.Host
	// 保存域名到变量
	if serverDomain == domain {
		return
	}
	serverDomain = domain

}

// Publish check token then save upload file to public directory
func Publish(token string, title string, data *multipart.FileHeader) (VideoListResponse, error) {
  username, err := utils.VerifyToken(token)
	if err != nil {
    videoListResponse := VideoListResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "token decoded error: " + err.Error()},
		}
		return videoListResponse, err
	}
	users, length, err := repository.NewUserDaoInstance().QueryUserByName(username)
	if length == 0 {
		videoListResponse := VideoListResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "User doesn't exist"},
		}
		return videoListResponse, err
	}
	//根据标题命名
	filename := filepath.Base(title)
	// 查询数据库视频当前主键
	id := repository.NewVideoDaoInstance().QueryVideoLatest()
	
  videoIdSequence = id
	//获取video序列
	// atomic.AddInt64(&videoIdSequence, 1)
	// user := usersLoginInfo[token]
	finalName := fmt.Sprintf("%d-%s.mp4", videoIdSequence, filename)
  _ = utils.WriteLog("likelog.txt", "filename: " + finalName)
	// 保存目录
	saveFile := filepath.Join("./public/videos", finalName)
  _ = utils.WriteLog("likelog.txt", "saveFile: " + saveFile)
	//视频url
	videourl := []string{"https:/", serverDomain, "static/video", finalName}
	playurl := strings.Join(videourl, "/")
  _ = utils.WriteLog("likelog.txt", "playurl: " + playurl)
	CoverName := fmt.Sprintf("%spng", finalName[0:len(finalName)-3])
	pngurl := []string{"https:/", serverDomain, "static/cover", CoverName}
	coverurl := strings.Join(pngurl, "/")
  _ = utils.WriteLog("likelog.txt", "coverurl: " + coverurl)
	newVideo := &repository.Video{
		Name:          users[0].Name,
		PlayUrl:       playurl,
		CoverUrl:      coverurl,
		FavoriteCount: 0,
		CommentCount:  0,
    IsFavorite:    false, 
		Title:         title,
		UploadTime:    time.Now(),
	}
	if err := repository.NewVideoDaoInstance().CreateVideo(newVideo); err != nil {
		videoListResponse := VideoListResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "insert video err:" + err.Error()},
		}
		return videoListResponse, err
	}
	srcFile, err := data.Open()
	if err != nil {
		videoListResponse := VideoListResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "open video err:" + err.Error()},
		}
		return videoListResponse, err
	}
	defer srcFile.Close()
	destFile, err := os.Create(saveFile)
	if err != nil {
		videoListResponse := VideoListResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "create file err:" + err.Error()},
		}
		return videoListResponse, err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		videoListResponse := VideoListResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "save video err:" + err.Error()},
		}
		return videoListResponse, err
	}
	// 抽帧做封面
	_, err = utils.GetSnapShot(saveFile, CoverName, 1)
	if err != nil {
		videoListResponse := VideoListResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "fail to abstract cover:" + err.Error()},
		}
		return videoListResponse, err
	}
	videoListResponse := VideoListResponse{
		Response: repository.Response{
			StatusCode: 0,
			StatusMsg:  finalName + " uploaded successfully"},
	}
	return videoListResponse, nil

}

// PublishList all users have same publish video list
func PublishList(token string) (VideoListResponse, error) {
  username, err := utils.VerifyToken(token)
	if err != nil {
    videoListResponse := VideoListResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "token decoded error: " + err.Error()},
		}
		return videoListResponse, err
	}

	_, _, err = repository.NewUserDaoInstance().QueryUserByName(username)
	if err != nil {
		videoListResponse := VideoListResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "query user err:" + err.Error()},
		}
		return videoListResponse, err
	}
	err = repository.NewVideoDaoInstance().UpdateVideoUrl(serverDomain)
	if err != nil {
		videoListResponse := VideoListResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "fail to update url!\n"},
		}
		return videoListResponse, err
	}
	videos, err := repository.NewVideoDaoInstance().QueryVideoByAuthor(username)
	if err != nil {
		videoListResponse := VideoListResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "query video list err:" + err.Error()},
		}
		return videoListResponse, err
	}
	VideoList, err := ConvertVideoDBToJSON(videos)
	if err != nil {
		videoListResponse := VideoListResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "videos are found, but Convert is failed"},
		}
		return videoListResponse, err
	}
	videoListResponse := VideoListResponse{
		Response: repository.Response{
			StatusCode: 0,
		},

		VideoList: VideoList,
		Token:     token,
	}
	return videoListResponse, nil
}
