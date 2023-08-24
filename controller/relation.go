package controller

import (
  "github.com/RaymondCode/simple-demo/repository"
  "github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	"net/http"
  "gorm.io/gorm"
)

type UserListResponse struct {
	repository.Response
	UserList []repository.User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")

	decoded_token, err := utils.VerifyToken(token)
  if err != nil{
    c.JSON(http.StatusOK, repository.Response{StatusCode: 1, StatusMsg: "Verify jwt error"})
    return 
  }
  _, _, err = repository.NewUserDaoInstance().QueryUserByName(decoded_token)
	if  err == gorm.ErrRecordNotFound{
		c.JSON(http.StatusOK, repository.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
		
  c.JSON(http.StatusOK, repository.Response{StatusCode: 0})
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: repository.Response{
			StatusCode: 0,
		},
		UserList: []repository.User{},
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: repository.Response{
			StatusCode: 0,
		},
		UserList: []repository.User{},
	})
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: repository.Response{
			StatusCode: 0,
		},
		UserList: []repository.User{},
	})
}
