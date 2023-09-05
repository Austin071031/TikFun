package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/repository"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
  "sync"
)


type VideoListResponse struct {
	repository.Response
	VideoList []utils.Video `json:"video_list"`
	Token     string  			`json:"token"`
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


var videoLock sync.Mutex 
// Publish check token then save upload file to public directory
func Publish(username string, title string, data *multipart.FileHeader) (VideoListResponse) {
	users, err := repository.NewUserDaoInstance().QueryUserByName(username)
	//根据标题命名
	filename := filepath.Base(title)
	// 查询数据库视频当前主键
	id := repository.NewVideoDaoInstance().QueryVideoLatest()
	
  videoIdSequence = id
	//获取video序列
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

  videoLock.Lock()
	if err := repository.NewVideoDaoInstance().CreateVideo(newVideo); err != nil {
		videoListResponse := VideoListResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "insert video err:" + err.Error()},
		}
		return videoListResponse
	}
	srcFile, err := data.Open()
	if err != nil {
		videoListResponse := VideoListResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "open video err:" + err.Error()},
		}
		return videoListResponse
	}
	defer srcFile.Close()
	destFile, err := os.Create(saveFile)
	if err != nil {
		videoListResponse := VideoListResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "create file err:" + err.Error()},
		}
		return videoListResponse
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, srcFile)
  
	if err != nil {
		videoListResponse := VideoListResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "save video err:" + err.Error()},
		}
		return videoListResponse
	}
	// 抽帧做封面
	_, err = utils.GetSnapShot(saveFile, CoverName, 1)
  videoLock.Unlock()
  
	if err != nil {
		videoListResponse := VideoListResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "fail to abstract cover:" + err.Error()},
		}
		return videoListResponse
	}
	videoListResponse := VideoListResponse{
		Response: repository.Response{
			StatusCode: 0,
			StatusMsg:  finalName + " uploaded successfully"},
	}
	return videoListResponse

}
var publishLock sync.Mutex 
// PublishList all users have same publish video list
func PublishList(username string) (VideoListResponse) {
	err := repository.NewVideoDaoInstance().UpdateVideoUrl(serverDomain)
	if err != nil {
		videoListResponse := VideoListResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "fail to update url!\n"},
		}
		return videoListResponse
	}
	videos, err := repository.NewVideoDaoInstance().QueryVideoByAuthor(username)
	if err != nil {
		videoListResponse := VideoListResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "query video list err:" + err.Error()},
		}
		return videoListResponse
	}

  publishLock.Lock()
  err = repository.NewUserDaoInstance().UpdateUserWorkCount(username, len((*videos)))
  publishLock.Unlock()
  if err != nil {
		videoListResponse := VideoListResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "update user work count err:" + err.Error()},
		}
		return videoListResponse
	}
  
	VideoList, err := utils.ConvertVideoDBToJSON(videos)
	if err != nil {
		videoListResponse := VideoListResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "videos are found, but Convert is failed"},
		}
		return videoListResponse
	}
	videoListResponse := VideoListResponse{
		Response: repository.Response{
			StatusCode: 0,
		},

		VideoList: VideoList,
	}
	return videoListResponse
}
