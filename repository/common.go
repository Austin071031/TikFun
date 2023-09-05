package repository
import "time"

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type LoginUserData struct {
	Username string `gorm:"column:name"`
  Password string `gorm:"column:password"`
	Token    string `gorm:"column:token"`
}

type Like struct {
	Name     string    `gorm:"column:name"`
	VideoId  int       `gorm:"column:videoId"`
	Liked    int64     `gorm:"column:liked"`
	Time     time.Time `gorm:"column:create_time"`
}

type User struct {
	Id              int64  `gorm:"column:id" json:"id"`
	Name            string `gorm:"column:name" json:"name"`
  Signature       string `gorm:"column:signature" json:"signature"`  
	FollowCount     int64  `gorm:"column:follow_count" json:"follow_count"`
	FollowerCount   int64  `gorm:"column:follower_count" json:"follower_count"`
	IsFollow        bool   `gorm:"column:isfollow" json:"is_follow"`
  Avatar          string `gorm:"column:avatar" json:"avatar"`
  BackgroundImage string `gorm:"column:background_image" json:"background_image"`
  FavoriteCount   int64  `gorm:"column:favorite_count" json:"favorite_count"` 
  TotalFavorited  string `gorm:"column:total_favorited" json:"total_favorited"`
  WorkCount       int64  `gorm:"column:work_count" json:"work_count"`
}

type Video struct {
	Id            int64     `gorm:"column:id" json:"id"`
	Name          string    `gorm:"column:name" json:"Name"`
	PlayUrl       string    `gorm:"column:playurl" json:"play_url`
	CoverUrl      string    `gorm:"column:coverurl" json:"cover_url"`
	FavoriteCount int64     `gorm:"column:favoritecount" json:"favorite_count"`
	CommentCount  int64     `gorm:"column:commentcount" json:"comment_count"`
  IsFavorite    bool      `gorm:"column:isfavorite" json:"is_favorite"`
	Title         string    `gorm:"column:title" json:"title"`
	UploadTime    time.Time `gorm:"column:uploadtime"`
}

type Comment struct {
	Id         int64  `gorm:"column:id" json:"id,omitempty"`
  UserId     int64  `gorm:"column:user_id" json:"-"`
  User       User   `gorm:"foreignkey:UserId;references:Id" json:"user"`
  VideoId    int64  `gorm:"column:video_id" json:"-"`
	Content    string `gorm:"column:content" json:"content,omitempty"`
	CreateDate string `gorm:"column:create_date" json:"create_date,omitempty"`
}

type Message struct {
	Id         int64  `json:"id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}

type MessageSendEvent struct {
	UserId     int64  `json:"user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type MessagePushEvent struct {
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}