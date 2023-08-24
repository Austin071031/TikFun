package repository

import (
	"sync"
  "gorm.io/gorm"
)


type CommentDao struct {
}

var (
	commentDao  *CommentDao
	commentOnce sync.Once
)

func NewCommentDaoInstance() *CommentDao {
	commentOnce.Do(func() {
		commentDao = &CommentDao{}
	})
	return commentDao
}

func (*CommentDao) CreateComment(comment *Comment) error {
	err := db.Table("comments").Create(struct {
		UserId     int64  `gorm:"column:user_id"`
		VideoId    int64  `gorm:"column:video_id"`
		Content    string `gorm:"column:content"`
		CreateDate string `gorm:"column:create_date"`
	}{comment.User.Id,comment.VideoId,comment.Content,comment.CreateDate}).Error
	if err != nil {
		return err
	}
	return nil
}

func (*CommentDao) DeleteCommentById(id int64) error {
	if err := db.Where("id = ?", id).Delete(&Comment{}).Error; err != nil {
		return err
	}
	return nil
}
//查询发布评论的用户id
func (*CommentDao) QueryUserIdById(id int64) (int64,error){
	var user_id int64
	if err := db.Table("comments").Select("user_id").Where("id=?", id).Find(&user_id).Error; err != nil {
		return 0, err
	}
	return user_id, nil
}
//查询video_id对应的所有的评论
func (*CommentDao) QueryCommentsByVideoId(video_id int64) ([]Comment, error) {
	var comments []Comment
	if err := db.Table("comments").Preload("User").Where("video_id=?",video_id).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

//查询video_id对应的评论数目
func (*CommentDao) QuerySumByVideoId(video_id int64) (int64, error) {
	var count int64
	if err := db.Table("comments").Select("sum(*)").Where("video_id=?",video_id).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (*CommentDao) UpdateCommentCountPlus(video_id int64) error{
  result := db.Table("videos").Where("id = ?", video_id).UpdateColumn("commentcount", gorm.Expr("commentcount  + ?", 1))

  if result.Error != nil {
    return result.Error
  }
  return nil
}

func (*CommentDao) UpdateCommentCountMinus(video_id int64) error{  
  result := db.Table("videos").Where("id = ?", video_id).UpdateColumn("commentcount", gorm.Expr("commentcount  - ?", 1))

  if result.Error != nil {
    return result.Error
  }
  return nil
}

// //查询发布视频的用户id
// func (*CommentDao) QueryAuthorIdByVideoId(video_id int64) (int64,error){
// 	var user_id int64
//   var name string
// 	if err := db.Table("videos").Select("name").Where("id=?", video_id).Find(&name).Error; err != nil {
// 		return 0, err
// 	}
//   if err := db.Table("users").Select("id").Where("name=?", name).Find(&user_id).Error; err != nil {
// 		return 0, err
// 	}
// 	return user_id, nil
// }