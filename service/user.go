package service

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/RaymondCode/simple-demo/repository"
	"github.com/RaymondCode/simple-demo/utils"
  "encoding/json"
)

type UserLoginResponse struct {
	repository.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	repository.Response
	User repository.User `json:"user"`
}

//用户注册
func Register(username, password string) (UserLoginResponse, error) {
  //查询用户是否存在
	_, length, err := repository.NewUserDaoInstance().QueryUserByName(username)
  //如果不存在，则创建用户
	if length == 0 {
		if inserterr := repository.NewUserDaoInstance().CreateUser(username, password); inserterr != nil {
			userloginresponse := UserLoginResponse{
				Response: repository.Response{StatusCode: 1, StatusMsg: "Insert user error" + err.Error()},
			}
			return userloginresponse, inserterr
		}
    //生成编码后的token，并储存
		SecretToken, _ := utils.GenerateToken(username)
		err := repository.NewUserDaoInstance().NewUserToken(username, SecretToken)
		if err != nil {
      userloginresponse := UserLoginResponse{
				Response: repository.Response{StatusCode: 1, StatusMsg: "Generate Token Error " + err.Error()},
      }
			return userloginresponse, err
		}
    //获取最新创建的用户
    last_user, err := repository.NewUserDaoInstance().QueryUserLast(username)
    if err != nil {
      userloginresponse := UserLoginResponse{
				Response: repository.Response{StatusCode: 1, StatusMsg: "Get last user Error " + err.Error()},
      }
			return userloginresponse, err
		}
    
    userloginresponse := UserLoginResponse{
			Response: repository.Response{StatusCode: 0},
			UserId:   last_user.Id,
			Token:    SecretToken,
		}
		return userloginresponse, nil
	} else {
    //如果存在，则返回用户已存在
		userloginresponse := UserLoginResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "User already exists"},
		}
		return userloginresponse, nil
	}
}

//用户登录
func Login(username, password string) (UserLoginResponse, error) {
	var PasswordIsTrue bool
  //查询该用户信息
	users, length, err := repository.NewUserDaoInstance().QueryUserByName(username)
	if length == 0 {
		userloginresponse := UserLoginResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "User doesn't exist",
			},
		}
		return userloginresponse, err
	}
  //用户输入的密码与数据库的密码做对比
	err = bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(password))
	if err == nil {
		PasswordIsTrue = true
	} else {
		PasswordIsTrue = false
	}
  //密码正确
	if PasswordIsTrue {
    //每次登录编码一次token以便传给url
		token, _ := utils.GenerateToken(username)
		err := repository.NewUserDaoInstance().UpdateUserToken(username, token)
		if err != nil {
			return UserLoginResponse{}, err
		}
		userloginresponse := UserLoginResponse{
			Response: repository.Response{
				StatusCode: 200,
			},
			Token:  token,
			UserId: users[0].Id,
		}
		return userloginresponse, nil
	} else {
    //密码不正确
		userloginresponse := UserLoginResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "Incorrect password",
			},
		}
		return userloginresponse, nil
	}
}

//返回用户信息，构建用户主页
func UserInfo(token string) (UserResponse, error) {
  //解码token获得username，以便查询此用户的信息。
  username, err := utils.VerifyToken(token)
  _ = utils.WriteLog("feed_querylike.txt", username)
  if err != nil {
    userresponse := UserResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "Token verified error"},
		}
		return userresponse, err
  }
  //根据用户名查询用户信息
	users, length, err := repository.NewUserDaoInstance().QueryUserByName(username)
	if length == 0 {
		userresponse := UserResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "User doesn't exist"},
		}
		return userresponse, err
	}
  jsondata, _ := json.Marshal(users[0])
  _ = utils.WriteLog("feed_querylike.txt", string(jsondata))

  
	userresponse := UserResponse{
		Response: repository.Response{StatusCode: 200},
		User:     users[0],
	}

  jsondata, _ = json.Marshal(userresponse)
  _ = utils.WriteLog("feed_querylike.txt", string(jsondata))
	return userresponse, nil
}
