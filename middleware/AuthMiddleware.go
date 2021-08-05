package middleware

import (
	"github.com/gin-gonic/gin"
	"godev/common"
	"godev/model"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc  {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg":"权限不足1"})
			c.Abort()
			return
		}

		tokenString = tokenString[7:]

		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"code":401, "msg":"权限不足2"})
			c.Abort()
			return
		}

		//token通过验证, 获取claims中的UserID
		userId := claims.UserId
		DB := common.GetDb()
		var user model.User
		DB.First(&user, userId)

		// 验证用户是否存在
		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "没有权限",
			})
			c.Abort()
			return
		}

		//用户存在 将user信息写入上下文
		c.Set("user", user)

		c.Next()

	}
}
