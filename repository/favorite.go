package repository

import (
	"sync"
  "time"
  "gorm.io/gorm"
)

func (Like) TableName() string { //Like结构体表示用户喜欢视频的记录
	return "likes"
}

type FavouriteDao struct { //FavouriteDao结构体表示用户喜欢或取消喜欢视频的数据访问层
}

var favouriteDao *FavouriteDao
var favouriteOnce sync.Once

func NewFavouriteInstance() *FavouriteDao {
	favouriteOnce.Do(
		func() {
			favouriteDao = &FavouriteDao{}
		})
	return favouriteDao
}

func (*FavouriteDao) CreateLiked(username string, videoId int) error {
	like := &Like{
		Name:     username,
		VideoId:  videoId,
		Liked:    1,
		Time:     time.Now(),
	}
	if err := db.Create(like).Error; err != nil {
		return err
	}
	return nil
}

func (*FavouriteDao) ChangeUnfavToFav(username string, videoId int) error {
	var like = 1
	result := db.Model(&Like{}).Where("name = ? AND VideoId = ?", username, videoId).Updates(map[string]interface{}{
		"liked": like,
		"Time":  time.Now(),
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (*FavouriteDao) ChangeFavToUnfav(username string, videoId int) error {
	var unlike = 0
	result := db.Model(&Like{}).Where("name = ? AND VideoId = ?", username, videoId).Updates(map[string]interface{}{
		"liked": unlike,
		"Time":  time.Now(),
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (*FavouriteDao) CheckLikedStatus(username string, videoId int) (bool, []Like) {
  var likevideo []Like
	//查询用户喜欢视频状态的数据访问层，检查用户和视频是否同时存在喜欢记录
	result := db.Where("name = ? AND VideoId = ?", username, videoId).Find(&likevideo) 
	//使用Find(&Like{})方法执行查询操作，并将结果赋值给result。
	return result.RowsAffected>0, likevideo  //通过查看result.RowsAffected的值来判断是否有记录受影响，如果大于0，则表示存在喜欢记录，返回true；
}



func (*FavouriteDao) GetFavoriteList(username string) ([] Video ,error) { //获取用户喜欢列表的数据访问层
  var videoList []Video
  err := db.Raw("SELECT v.id as Id, v.name as Name, v.playurl as PlayUrl, v.coverurl as CoverUrl, v.favoritecount as FavoriteCount, v.commentcount as CommentCount,v.isfavorite as IsFavorite,v.title as Title,v.uploadtime as UploadTime FROM likes l JOIN videos v ON l.videoId = v.id WHERE l.liked = ? ORDER BY l.create_time DESC",
		1).Scan(&videoList)
  
	// err := db.Model(&Like{}).Where("Username = ? AND liked = ?", username, 1).Select("videoId").Find(&videoIds).Error
	//通过Where方法指定查询条件，查询用户喜欢的视频的videoId字段，通过Select方法指定查询结果只返回videoId字段，并使用Find(&videoIds)方法执行查询操作，将结果保存到videoIds变量中。
	if err != nil {
		// util.Logger.Error("get favorite list err: " + err.Error())
		return nil, err.Error
	} //如果查询过程中出现错误，则记录错误信息并返回错误；否则，返回videoIds表示查询到的视频ID列表。
	return videoList, nil
}

func (*FavouriteDao) UpdateFavouriteCountPlus(videoId int) error{
  result := db.Model(&Video{}).Where("id = ?", videoId).UpdateColumn("favoritecount", gorm.Expr("favoritecount  + ?", 1))

  if result.Error != nil {
    return result.Error
  }
  return nil
}

func (*FavouriteDao) UpdateFavouriteCountMinus(videoId int) error{  
  result := db.Model(&Video{}).Where("id = ?", videoId).UpdateColumn("favoritecount", gorm.Expr("favoritecount  - ?", 1))

  if result.Error != nil {
    return result.Error
  }
  return nil
}


func(*FavouriteDao) UpdateUserLikedVideo(username string)(int64, error){
  var count int64
  result := db.Where("name = ? and liked = ?", username, 1).Count(&count)
  if result.Error != nil {
    return 0, result.Error
  }
  return count, nil
}


func(*FavouriteDao) FindUserLikedVideo(username string)([]int, error){
  var LikeVideoId []int
  result := db.Model(Like{}).Where("name=? and liked = ?",username,1).Select("videoId").Find(&LikeVideoId)
  if result.Error != nil {
    return LikeVideoId, result.Error
  }
  return LikeVideoId,nil
}


  
