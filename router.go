package main

import (
  "github.com/RaymondCode/simple-demo/controller"
  "github.com/RaymondCode/simple-demo/service"
  "github.com/RaymondCode/simple-demo/middleware/jwt"
  "github.com/gin-gonic/gin"
  "net/http"
)

func initRouter(r *gin.Engine) {
  // public directory is used to serve static resources
 	r.Static("/static/video", "./public/videos")
	r.Static("/static/cover", "./public/covers")
  r.LoadHTMLGlob("templates/*")

  // home page
  r.GET("/", func(c *gin.Context) {
    c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
  })
  
  apiRouter := r.Group("/douyin")
  {
    // basic apis
    apiRouter.GET("/feed/",service.GetServerDomain, controller.Feed)
    apiRouter.GET("/user/", controller.UserInfo)
    apiRouter.POST("/user/register/", controller.Register)
    apiRouter.POST("/user/login/", controller.Login)

    publishRouter := apiRouter.Group("/publish")
    {
      publishRouter.Use(jwt.JWT_PUBLISH())
      publishRouter.POST("/action/", controller.Publish)
      publishRouter.GET("/list/", service.GetServerDomain, controller.PublishList)
    }
  
    // extra apis - I
    favoriteRouter := apiRouter.Group("/favorite")
    {
      favoriteRouter.Use(jwt.JWT())
      favoriteRouter.POST("/action/", controller.FavoriteAction)
      favoriteRouter.GET("/list/", controller.FavoriteList)
    }

    commentRouter := apiRouter.Group("/comment")
    {
      commentRouter.Use(jwt.JWT())
      commentRouter.POST("/action/", controller.CommentAction)
    }
    apiRouter.GET("/comment/list/", controller.CommentList)
  
    // extra apis - II
    apiRouter.POST("/relation/action/", controller.RelationAction)
    apiRouter.GET("/relation/follow/list/", controller.FollowList)
    apiRouter.GET("/relation/follower/list/", controller.FollowerList)
    apiRouter.GET("/relation/friend/list/", controller.FriendList)
    apiRouter.GET("/message/chat/", controller.MessageChat)
    apiRouter.POST("/message/action/", controller.MessageAction) 
  }
}
