package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"godev/response"
)

func RecoveryMiddleware() gin.HandlerFunc  {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				response.Fail(ctx, fmt.Sprint(err), nil)
			}
		}()
		ctx.Next()
	}
}
