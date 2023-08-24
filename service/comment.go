package service
import(
	"strconv"
	"github.com/RaymondCode/simple-demo/repository"
	"time"
	"fmt"
)

type CommentListResponse struct {
	repository.Response
	CommentList []repository.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	repository.Response
	Comment repository.Comment `json:"comment,omitempty"`
}

func CreateComment(content string,video_id_text string,user *repository.User)(CommentActionResponse,error){
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
		User:       *user,
		VideoId:    video_id,
		Content:    content,
		CreateDate: time.Now().Format("01-02"),
	}
	commentDao:=repository.NewCommentDaoInstance()
	if err=commentDao.CreateComment(comment);err!=nil{
		return CommentActionResponse{
			Response:    repository.Response{
				StatusCode: 1,
				StatusMsg:  "Comment creation failed",
			},
		},err
	}
  //Video结构体CommentCount+1
  
  commentDao.UpdateCommentCountPlus(video_id)
  
	return CommentActionResponse{
		Response: repository.Response{StatusCode: 0},
		Comment: *comment,
	},nil
}


func DeleteComment(comment_id_text string,video_id_text string,user *repository.User) (CommentActionResponse,error) {
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
	commentDao:=repository.NewCommentDaoInstance()
	comment_user_id,err:=commentDao.QueryUserIdById(comment_id)
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
	if comment_user_id==user.Id {//评论发布者
		err = commentDao.DeleteCommentById(comment_id)
		if err != nil {
			return CommentActionResponse{
				Response:    repository.Response{
					StatusCode: 1,
					StatusMsg:  "comment deletion failed",
				},
			},err
		}
    commentDao.UpdateCommentCountMinus(video_id)
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
	commentDao:=repository.NewCommentDaoInstance()
	comments,err:=commentDao.QueryCommentsByVideoId(video_id)
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