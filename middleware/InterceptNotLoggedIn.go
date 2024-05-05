package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rookiefront/api-core/service"
)

func InterceptNotLoggedIn(c *gin.Context) {
	token := c.GetHeader("X-Token")
	userMap, err := service.User.ParseToken(token)
	if err != nil {
		c.Abort()
		c.JSON(200, gin.H{
			"msg":  "未登录",
			"code": 400001,
			"data": nil,
		})
		return
	}
	c.Set("uid", fmt.Sprintf("%.0f", userMap["id"]))
	c.Next()
}
