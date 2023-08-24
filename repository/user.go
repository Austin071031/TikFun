package repository

import (
	"golang.org/x/crypto/bcrypt"
	"sync"
	// "time"
	_ "encoding/base64"
)

func (User) TableName() string {
	return "users"
}

func (LoginUserData) TableName() string {
	return "loginuserdata"
}

type UserDao struct {
}

var userDao *UserDao
var userOnce sync.Once

func NewUserDaoInstance() *UserDao {
	userOnce.Do(
		func() {
			userDao = &UserDao{}
		})
	return userDao
}


//根据用户名字查找用户
func (*UserDao) QueryUserByName(name string) ([]User, int64, error) {
	var users []User
	result := db.Where("name = ?", name).Find(&users)
	if result.Error != nil {
		// util.Logger.Error("find user by id err:" + err.Error())
		return nil, 0, result.Error
	}
	return users, result.RowsAffected, nil
}

//注册，新建一个用户
func (*UserDao) CreateUser(username, password string) error {
	Secretpassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	//encodedPassword := base64.StdEncoding.EncodeToString(Secretpassword)
	user := &User{
		Name:          username,
		Password:      string(Secretpassword),
    Signature:     "",
		FollowCount:   0,
		FollowerCount: 0,
    IsFollow:      false,
    Avatar:        "",
    BackgroundImage:"",
    FavoriteCount: 0,
    TotalFavorited:"",
    WorkCount:0,
    
	}

	if err := db.Create(user).Error; err != nil {
		// util.Logger.Error("insert post err:" + err.Error())
		return err
	}
	return nil
}


//用户注册后，给他一个token
func (*UserDao) NewUserToken(username, token string) error {
	loginuserdata := &LoginUserData{
		Username: username,
		Token:    token,
	}
	if err := db.Create(loginuserdata).Error; err != nil {
		// util.Logger.Error("insert post err:" + err.Error())
		return err
	}
	return nil
}

//用户每次登录更新token
func (*UserDao) UpdateUserToken(username, token string) error {
	// 定义更新字段和新值的映射
	updates := map[string]interface{}{
		"Token": token,
	}
	// 执行更新操作，将 Token 更新为新值
	if err := db.Model(&LoginUserData{}).Where("name = ?", username).Updates(updates).Error; err != nil {
		// 处理更新错误
		return err
	}
	return nil
}


func (*UserDao) QueryUserByToken(token string) (*User, error) {
	var user User
	err := db.Where("token = ?", token).First(&user).Error
	if err != nil {
		// util.Logger.Error("find user by id err:" + err.Error())
		return nil, err
	}
	return &user, nil
}

func (*UserDao) QueryUserLast(username string) (*User, error) {
	var user User
	result := db.Last(&user)
	if result.Error != nil {
		// util.Logger.Error("find user by id err:" + err.Error())
		return nil, result.Error
	}
	return &user, nil
}



