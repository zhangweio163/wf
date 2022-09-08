package http

import (
	"Reginster/server"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAll() map[string][]*server.Server {
	c := server.NewRegServer()
	return c.ServerList
}
func HttpServer(port int) {
	r := gin.Default()
	r.GET("/user", func(c *gin.Context) {
		c.JSON(http.StatusOK, GetAll())
	})
	r.Run(fmt.Sprintf(":%d", port))
}
