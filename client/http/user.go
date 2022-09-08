package http

import (
	"Userclient/endport"
	"Userclient/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func RunServer(port int) {
	r := gin.Default()
	r.GET("/user", func(c *gin.Context) {
		var reqUser *user.UserReq
		userId := c.Query("id")
		id, _ := strconv.Atoi(userId)
		reqUser = &user.UserReq{Userid: int32(id)}
		c.JSON(http.StatusOK, endport.GetUserEndPort(reqUser))
	})
	r.Run(fmt.Sprintf(":%d", port))
}
