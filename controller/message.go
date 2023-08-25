package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
  "github.com/RaymondCode/simple-demo/repository"
  "github.com/RaymondCode/simple-demo/utils"
  "gorm.io/gorm"
)

var tempChat = map[string][]repository.Message{}

var messageIdSequence = int64(1)

type ChatResponse struct {
	repository.Response
	MessageList []repository.Message `json:"message_list"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	content := c.Query("content")

	username, err := utils.VerifyToken(token)
  if err != nil{
    c.JSON(http.StatusOK, repository.Response{StatusCode: 1, StatusMsg: "Verify jwt error"})
    return 
  }
  users, err := repository.NewUserDaoInstance().QueryUserByName(username)
	if  len(users) == 0{
		c.JSON(http.StatusOK, repository.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
  
  userIdB, _ := strconv.Atoi(toUserId)
  chatKey := genChatKey(users[0].Id, int64(userIdB))

  atomic.AddInt64(&messageIdSequence, 1)
  curMessage := repository.Message{
    Id:         messageIdSequence,
    Content:    content,
    CreateTime: time.Now().Format(time.Kitchen),
  }

  if messages, exist := tempChat[chatKey]; exist {
    tempChat[chatKey] = append(messages, curMessage)
  } else {
    tempChat[chatKey] = []repository.Message{curMessage}
  }
  c.JSON(http.StatusOK, repository.Response{StatusCode: 0})
}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")

	if user, err := repository.NewUserDaoInstance().QueryUserByToken(token); err != gorm.ErrRecordNotFound{
		userIdB, _ := strconv.Atoi(toUserId)
		chatKey := genChatKey(user.Id, int64(userIdB))

		c.JSON(http.StatusOK, ChatResponse{Response: repository.Response{StatusCode: 0}, MessageList: tempChat[chatKey]})
	} else {
		c.JSON(http.StatusOK, repository.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

func genChatKey(userIdA int64, userIdB int64) string {
	if userIdA > userIdB {
		return fmt.Sprintf("%d_%d", userIdB, userIdA)
	}
	return fmt.Sprintf("%d_%d", userIdA, userIdB)
}
