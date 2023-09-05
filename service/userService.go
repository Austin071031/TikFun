package service

import (
	"sync"
	"github.com/RaymondCode/simple-demo/repository"
	"github.com/RaymondCode/simple-demo/utils"
	"golang.org/x/crypto/bcrypt"
	// "encoding/json"
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

var registerLock sync.Mutex

//用户注册
func Register(username, password string) (UserLoginResponse) {
  //查询用户是否存在
	users, err := repository.NewUserDaoInstance().QueryUserByName(username)
	if err != nil{
		userloginresponse := UserLoginResponse{
				Response: repository.Response{StatusCode: 1, StatusMsg: "Query user exists error." + err.Error()},
			}
			return userloginresponse
	}
  //如果不存在，则创建用户
	if len(users) == 0 {
    //用户注册并发处理
    registerLock.Lock()
    newUser, err := repository.NewUserDaoInstance().CreateUser(username)
    registerLock.Unlock()
    
		if err != nil {
			userloginresponse := UserLoginResponse{
				Response: repository.Response{StatusCode: 1, StatusMsg: "Insert user error." + err.Error()},
			}
			return userloginresponse
		}
    //生成编码后的token，并储存
    registerLock.Lock()
		SecretToken, _ := utils.GenerateToken(username)
		err = repository.NewUserDaoInstance().CreateUserToken(username, password, SecretToken)
    registerLock.Unlock()
		if err != nil {
      userloginresponse := UserLoginResponse{
				Response: repository.Response{StatusCode: 1, StatusMsg: "Generate token error." + err.Error()},
      }
			return userloginresponse
		}
    
    userloginresponse := UserLoginResponse{
			Response: repository.Response{StatusCode: 0},
			UserId:   newUser.Id,
			Token:    SecretToken,
		}
		return userloginresponse
	} else {
    //如果存在，则返回用户已存在
		userloginresponse := UserLoginResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "User already exists."},
		}
		return userloginresponse
	}
}

var loginLock sync.Mutex
//用户登录
func Login(username, password string) (UserLoginResponse) {
	var PasswordIsTrue bool
  //查询该用户信息
	usersdata, length, err := repository.NewUserDaoInstance().QueryUserDataByName(username)
	if length == 0 {
		userloginresponse := UserLoginResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "User doesn't exist.",
			},
		}
		return userloginresponse
	}
  //用户输入的密码与数据库的密码做对比
	err = bcrypt.CompareHashAndPassword([]byte(usersdata[0].Password), []byte(password))
	if err == nil {
		PasswordIsTrue = true
	} else {
		PasswordIsTrue = false
	}
  //密码正确
	if PasswordIsTrue {
    //每次登录编码一次token以便传给url
    loginLock.Lock()
		token, _ := utils.GenerateToken(username)
		err := repository.NewUserDaoInstance().UpdateUserToken(username, token)
    loginLock.Unlock()
    
		if err != nil {
			userloginresponse := UserLoginResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "Login update user token error." + err.Error(),
			},
		}
		return userloginresponse
		}

    users, err := repository.NewUserDaoInstance().QueryUserByName(username)
  	if err != nil {
  		userloginresponse := UserLoginResponse{
  			Response: repository.Response{
  				StatusCode: 1,
  				StatusMsg:  "Get user info error." + err.Error(),
  			},
  		}
  		return userloginresponse
  	}
		userloginresponse := UserLoginResponse{
			Response: repository.Response{
				StatusCode: 0,
			},
			Token:  token,
			UserId: users[0].Id,
		}
    
		return userloginresponse
	} else {
    //密码不正确
		userloginresponse := UserLoginResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "Incorrect password.",
			},
		}
		return userloginresponse
	}
}

//返回用户信息，构建用户主页
func UserInfo(token string) (UserResponse) {
  //解码token获得username，以便查询此用户的信息。
  username, err := utils.VerifyToken(token)

  if err != nil {
    userresponse := UserResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "Token verified erro." + err.Error()},
		}
		return userresponse
  }
  //根据用户名查询用户信息
	users, err := repository.NewUserDaoInstance().QueryUserByName(username)
	if len(users) == 0 {
		userresponse := UserResponse{
			Response: repository.Response{
				StatusCode: 1,
				StatusMsg:  "User doesn't exist." + err.Error()},
		}
		return userresponse
	}
  
	userresponse := UserResponse{
		Response: repository.Response{StatusCode: 0, StatusMsg:  "Get user info success"},
		User:     users[0],
	}
  
	return userresponse
}
