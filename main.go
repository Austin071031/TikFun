package main

import (
  "github.com/RaymondCode/simple-demo/repository"
	"github.com/RaymondCode/simple-demo/service"
  //"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
  "os"
)

func main() {
  if err := repository.Init(); err != nil {
		os.Exit(-1)
	}
  go service.RunMessageServer()

  r := gin.Default()

  initRouter(r)

  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
