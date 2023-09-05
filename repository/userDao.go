package repository

import (
	"golang.org/x/crypto/bcrypt"
	"sync"
	// "time"
	_ "encoding/base64"
  "strconv"
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


//根据用户名字查找用户登录信息
func (*UserDao) QueryUserDataByName(name string) ([]LoginUserData, int64, error) {
	var usersdata []LoginUserData
	result := db.Where("name = ?", name).Find(&usersdata)
	if result.Error != nil {
		// util.Logger.Error("find user by id err:" + err.Error())
		return nil, 0, result.Error
	}
	return usersdata, result.RowsAffected, nil
}

//根据用户名字查找用户
func (*UserDao) QueryUserByName(name string) ([]User, error) {
	var users []User
	result := db.Where("name = ?", name).Find(&users)
	if result.Error != nil {
		// util.Logger.Error("find user by id err:" + err.Error())
		return nil, result.Error
	}
	return users, nil
}

//注册，新建一个用户
func (*UserDao) CreateUser(username string) (User, error) {
	//encodedPassword := base64.StdEncoding.EncodeToString(Secretpassword)
	user := User{
		Name:          username,
    Signature:     "",
		FollowCount:   0,
		FollowerCount: 0,
    IsFollow:      false,
    Avatar:        "",
    BackgroundImage:"",
    FavoriteCount: 0,
    TotalFavorited:"0",
    WorkCount:0,
    
	}

	if err := db.Create(&user).Error; err != nil {
		// util.Logger.Error("insert post err:" + err.Error())
		return User{}, err
	}
	return user, nil
}


//用户注册后，给他一个token
func (*UserDao) CreateUserToken(username, password, token string) error {
	Secretpassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  loginuserdata := &LoginUserData{
		Username: username,
    Password: string(Secretpassword),
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


func (*UserDao) UpdateUserFavoriteCount(username string, favoriteCount int) (error) {
  updates := map[string]interface{}{
		"FavoriteCount": favoriteCount,
	}
	if err := db.Model(&User{}).Where("name = ?", username).Updates(updates).Error; err != nil {
		// 处理更新错误
		return err
	}
	return nil
}

func (*UserDao) UpdateUserTotalFavorited(authorname string, flag int) (error) {
  var totfavorited []string
  result := db.Model(&User{}).Where("name=?", authorname).Select("total_favorited").Find(&totfavorited)
  if result.Error != nil {
		// util.Logger.Error("find user by id err:" + err.Error())
		return result.Error
	}

  var totf string
  totf = totfavorited[0]
  temp, err:= strconv.Atoi(totf)
  if err != nil{
    return err
  }

  updates := map[string]interface{}{
    "TotalFavorited": "0",
  }
  if flag == 0{
		updates["TotalFavorited"] = strconv.Itoa(temp - 1)
  }else{
    updates["TotalFavorited"] = strconv.Itoa(temp + 1)
  }
  
  

  if err := db.Model(&User{}).Where("name = ?", authorname).Updates(updates).Error; err != nil {
		// 处理更新错误
		return err
	}
	return nil
}

func (*UserDao) UpdateUserWorkCount(username string, workCount int) (error) {
  updates := map[string]interface{}{
		"WorkCount": workCount,
	}
	if err := db.Model(&User{}).Where("name = ?", username).Updates(updates).Error; err != nil {
		// 处理更新错误
		return err
	}
	return nil
}

  
  

