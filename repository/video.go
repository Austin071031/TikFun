package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"sync"
  "database/sql"
  "strconv"

)

func (Video) TableName() string {
	return "videos"
}

type VideoDao struct {
}

var videoDao *VideoDao
var videoOnce sync.Once

func NewVideoDaoInstance() *VideoDao {
	videoOnce.Do(
		func() {
			videoDao = &VideoDao{}
		})
	return videoDao
}
func (*VideoDao) QueryVideoByAuthor(username string) (*[]Video, error) {
	var videoList []Video
	err := db.Where("name = ?", username).Find(&videoList).Error
	if errors.Is(gorm.ErrRecordNotFound, err) {
		return nil, err
	}
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &videoList, err
}

// QueryVideoFeed 按照Feed要求返回视频列表
func (*VideoDao) QueryVideoFeed() (*[]Video, error) {
	var videoList []Video
	err := db.Order("uploadtime desc").Limit(30).Find(&videoList).Error
	if errors.Is(gorm.ErrRecordNotFound, err) {
		return nil, err
	}
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &videoList, err
}

// QueryVideoLatest 查询下一个主键ID
func (*VideoDao) QueryVideoLatest() (id int64) {
  //sql.NullInt64
  var latest_id sql.NullInt64   
	_ = db.Table("videos").Select("MAX(id)").Find(&latest_id)
  if !latest_id.Valid{
    return 1
  }else{
    return latest_id.Int64 + 1
  }
}

// CreateVideo 在数据库中创建记录
func (*VideoDao) CreateVideo(video *Video) error {
	if err := db.Create(video).Error; err != nil {
		return err
	}
	return nil
}

// UpdateVideoUrl 更新数据库中的URL
func (*VideoDao) UpdateVideoUrl(domain string) error {
	var videos []Video
	if err := db.Find(&videos).Error; err != nil {
		return err
	}
	for _, video := range videos {
		fileName := fmt.Sprintf("%d-%s.mp4", video.Id, video.Title)
		videourl := []string{"https:/", domain, "static/video", fileName}
		playurl := strings.Join(videourl, "/")
		CoverName := fmt.Sprintf("%spng", fileName[0:len(fileName)-3])
		pngurl := []string{"https:/", domain, "static/cover", CoverName}
		coverurl := strings.Join(pngurl, "/")
		updates := map[string]interface{}{
			"playurl":  playurl,
			"coverurl": coverurl,
		}
		if err := db.Model(&Video{}).Where("id = ?", video.Id).Updates(updates).Error; err != nil {
			return err
		}
	}
	return nil
}

func(*VideoDao) UpdateUserFavcount(username string)(string, int64, error){
  var videos []Video
  // 查询用户名下所有视频
  result := db.Where("name = ?", username).Find(&videos)
  if result.Error != nil {
    return "",0, result.Error
  }
  // 点赞总数
  var totalFavCount int64 = 0
  //作品总数
  var videocount int64 = 0
  // 遍历视频,累加点赞数
  for _, video := range videos {
    totalFavCount += video.FavoriteCount 
    videocount ++
  }
  str := strconv.FormatInt(totalFavCount, 10)
  
  return str,videocount, nil
}



