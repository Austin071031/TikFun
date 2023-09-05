package service

import (
	"strconv"
	"time"
	"github.com/RaymondCode/simple-demo/repository"
  "sync"
)

type CommentListResponse struct {
	repository.Response
	CommentList []repository.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	repository.Response
	Comment repository.Comment `json:"comment,omitempty"`
}

var commentLock sync.Mutex 

func CreateComment(content string,video_id_text string,username string)(CommentActionResponse){
  users, err := repository.NewUserDaoInstance().QueryUserByName(username)
	if len(users) == 0{
		return CommentActionResponse{
			Response:    repository.Response{
				StatusCode: 1,
				StatusMsg:  "find comment author failed",
			},
		}
	}
  
	video_id,err:=strconv.ParseInt(video_id_text, 10, 64)
	if err != nil {
		return CommentActionResponse{
			Response:    repository.Response{
				StatusCode: 1,
				StatusMsg:  "Video_id read failed",
			},
		}
	}
	comment:=&repository.Comment{
		User:       users[0],
		VideoId:    video_id,
		Content:    content,
		CreateDate: time.Now().Format("01-02"),
	}
	// commentDao:=repository.NewCommentDaoInstance()
  commentLock.Lock()
	if err=repository.NewCommentDaoInstance().CreateComment(comment);err!=nil{
		return CommentActionResponse{
			Response:    repository.Response{
				StatusCode: 1,
				StatusMsg:  "Comment creation failed",
			},
		}
	}
  
  //Video结构体CommentCount+1
  repository.NewCommentDaoInstance().UpdateCommentCountPlus(video_id)
  commentLock.Unlock()
	return CommentActionResponse{
		Response: repository.Response{StatusCode: 0},
		Comment: *comment,
	}
}


func DeleteComment(comment_id_text string,video_id_text string,username string) (CommentActionResponse) {
  users, err := repository.NewUserDaoInstance().QueryUserByName(username)
	if len(users) == 0{
		return CommentActionResponse{
			Response:    repository.Response{
				StatusCode: 1,
				StatusMsg:  "find comment author failed",
			},
		}
	}
  
	comment_id,err:=strconv.ParseInt(comment_id_text, 10, 64)
	if err != nil {
		return CommentActionResponse{
			Response:    repository.Response{
				StatusCode: 1,
				StatusMsg:  "comment_id read failed",
			},
		}
	}
  video_id,err:=strconv.ParseInt(video_id_text, 10, 64)
	if err != nil {
		return CommentActionResponse{
			Response:    repository.Response{
				StatusCode: 1,
				StatusMsg:  "Video_id read failed",
			},
		}
	}

	comment_user_id,err:=repository.NewCommentDaoInstance().QueryUserIdById(comment_id)
	if err != nil {
		return CommentActionResponse{
			Response:    repository.Response{
				StatusCode: 1,
				StatusMsg:  "owner of comment read failed",
			},
		}
	}
  
	if comment_user_id==users[0].Id {//评论发布者
    commentLock.Lock()
		err = repository.NewCommentDaoInstance().DeleteCommentById(comment_id)
		if err != nil {
			return CommentActionResponse{
				Response:    repository.Response{
					StatusCode: 1,
					StatusMsg:  "comment deletion failed",
				},
			}
		}
    repository.NewCommentDaoInstance().UpdateCommentCountMinus(video_id)
    commentLock.Unlock()
    
		return CommentActionResponse{
      Response:   repository.Response{
			StatusCode: 0,
      },
    }
      
	}
  //其他人&&comment_id=0
	return CommentActionResponse{
		Response:    repository.Response{
			StatusCode: 1,
			StatusMsg:  "Error",
		},
	}
}

func CommentList(video_id_text string) (CommentListResponse){
	video_id,err:=strconv.ParseInt(video_id_text, 10, 64)
	if err != nil {
		return CommentListResponse{
			Response:    repository.Response{
				StatusCode: 1,
				StatusMsg:  "Video_id read failed",
			},
		}
	}
	comments,err:=repository.NewCommentDaoInstance().QueryCommentsByVideoId(video_id)
	if err != nil{
		return CommentListResponse{
			Response:    repository.Response{
				StatusCode: 1,
				StatusMsg:  "comments query failed",
			},
		}
	}
	return CommentListResponse{
		Response:    repository.Response{StatusCode: 0},
		CommentList: comments,
	}
}