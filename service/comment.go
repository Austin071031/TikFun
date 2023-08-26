package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/RaymondCode/simple-demo/repository"
	"github.com/RaymondCode/simple-demo/utils"
)

type CommentListResponse struct {
	repository.Response
	CommentList []repository.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	repository.Response
	Comment repository.Comment `json:"comment,omitempty"`
}

func CreateComment(content string,video_id_text string,username string)(CommentActionResponse,error){
  users, err := repository.NewUserDaoInstance().QueryUserByName(username)
	if len(users) == 0{
		return CommentActionResponse{
			Response:    repository.Response{
				StatusCode: 1,
				StatusMsg:  "find comment author failed",
			},
		},err
	}
  
	video_id,err:=strconv.ParseInt(video_id_text, 10, 64)
	if err != nil {
		return CommentActionResponse{
			Response:    repository.Response{
				StatusCode: 1,
				StatusMsg:  "Video_id read failed",
			},
		},err
	}
	comment:=&repository.Comment{
		User:       users[0],
		VideoId:    video_id,
		Content:    content,
		CreateDate: time.Now().Format("01-02"),
	}
	// commentDao:=repository.NewCommentDaoInstance()
	if err=repository.NewCommentDaoInstance().CreateComment(comment);err!=nil{
		return CommentActionResponse{
			Response:    repository.Response{
				StatusCode: 1,
				StatusMsg:  "Comment creation failed",
			},
		},err
	}
  //Video结构体CommentCount+1
  
  repository.NewCommentDaoInstance().UpdateCommentCountPlus(video_id)
  
	return CommentActionResponse{
		Response: repository.Response{StatusCode: 0},
		Comment: *comment,
	},nil
}


func DeleteComment(comment_id_text string,video_id_text string,username string) (CommentActionResponse,error) {
  users, err := repository.NewUserDaoInstance().QueryUserByName(username)
	if len(users) == 0{
		return CommentActionResponse{
			Response:    repository.Response{
				StatusCode: 1,
				StatusMsg:  "find comment author failed",
			},
		},err
	}
  
	//判断用户是视频发布者/评论发布者/其他人
	comment_id,err:=strconv.ParseInt(comment_id_text, 10, 64)
	if err != nil {
		return CommentActionResponse{
			Response:    repository.Response{
				StatusCode: 1,
				StatusMsg:  "comment_id read failed",
			},
		},err
	}
  video_id,err:=strconv.ParseInt(video_id_text, 10, 64)
	if err != nil {
		return CommentActionResponse{
			Response:    repository.Response{
				StatusCode: 1,
				StatusMsg:  "Video_id read failed",
			},
		},err
	}
	//评论的发布者
	// commentDao:=repository.NewCommentDaoInstance()
  
	comment_user_id,err:=repository.NewCommentDaoInstance().QueryUserIdById(comment_id)
	if err != nil {
		return CommentActionResponse{
			Response:    repository.Response{
				StatusCode: 1,
				StatusMsg:  "owner of comment read failed",
			},
		},err
	}
  //视频的发布者
 //  video_user_id,err:=commentDao.QueryAuthorIdByVideoId(video_id)
	// if err != nil {
	// 	return CommentActionResponse{
	// 		Response:    repository.Response{
	// 			StatusCode: 1,
	// 			StatusMsg:  "owner of video read failed",
	// 		},
	// 	},err
	// }
  _ = utils.WriteLog("commentaction.txt", strconv.FormatInt(comment_user_id,10))
  _ = utils.WriteLog("commentaction.txt", strconv.FormatInt(users[0].Id,10))
  
	if comment_user_id==users[0].Id {//评论发布者
		err = repository.NewCommentDaoInstance().DeleteCommentById(comment_id)
		if err != nil {
			return CommentActionResponse{
				Response:    repository.Response{
					StatusCode: 1,
					StatusMsg:  "comment deletion failed",
				},
			},err
		}
    repository.NewCommentDaoInstance().UpdateCommentCountMinus(video_id)
		return CommentActionResponse{},nil
	}
  //其他人
	return CommentActionResponse{
		Response:    repository.Response{
			StatusCode: 1,
			StatusMsg:  "Unauthorized user",
		},
	},fmt.Errorf("Unauthorized user")
}

func CommentList(video_id_text string) (CommentListResponse,error){
	video_id,err:=strconv.ParseInt(video_id_text, 10, 64)
	if err != nil {
		return CommentListResponse{
			Response:    repository.Response{
				StatusCode: 1,
				StatusMsg:  "Video_id read failed",
			},
		},err
	}
	// commentDao:=repository.NewCommentDaoInstance()
	comments,err:=repository.NewCommentDaoInstance().QueryCommentsByVideoId(video_id)
	if err != nil{
		return CommentListResponse{
			Response:    repository.Response{
				StatusCode: 1,
				StatusMsg:  "comments query failed",
			},
		},err
	}
	return CommentListResponse{
		Response:    repository.Response{StatusCode: 0},
		CommentList: comments,
	},nil
}