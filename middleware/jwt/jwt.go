package jwt

import(
  "github.com/RaymondCode/simple-demo/utils"
  "github.com/RaymondCode/simple-demo/repository"
  "net/http"
  "github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc{
  return func(c *gin.Context){
    token := c.Query("token")
    
    if token == ""{
      c.JSON(http.StatusUnauthorized, gin.H{
        "code": 1,
        "msg": "User unauthorized",
      })
      c.Abort()
      return
    }else{
      username, err := utils.VerifyToken(token)
      if err != nil{
        c.JSON(http.StatusOK, gin.H{
          "code": 1,
          "msg": "Verified user token error",
        })
        c.Abort()
        return
      }
      users, err := repository.NewUserDaoInstance().QueryUserByName(username)
      if  len(users) == 0{
        c.JSON(http.StatusOK, gin.H{
          "code": 1,
          "msg": "User doesn't exist",
        })
        c.Abort()
        return
      }
      c.Set("username", username)
    }
    c.Next()
  }
}