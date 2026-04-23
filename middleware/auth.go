package middleware

import (
	"fmt"
	"gin-user/pkg/e"
	jwt2 "gin-user/pkg/jwt"
	"gin-user/pkg/utils/ctl"
	"github.com/gin-gonic/gin"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		var code int
		code = e.SUCCESS
		accessToken := extraceToken(c)

		fmt.Printf("access_token:%v\n", accessToken)

		if accessToken == "" {
			code = e.InvalidParams
			c.JSON(200, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   "Token不能为空",
			})
			c.Abort()
			return
		}

		jwt := jwt2.NewJwtManager()
		claims, err := jwt.ParseToken(accessToken)

		if err != nil {
			code = e.ErrorAuthCheckTokenFail
			c.JSON(200, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   "鉴权失败",
				"error":  err.Error(),
			})
			c.Abort()
			return
		}

		c.Request = c.Request.WithContext(ctl.NewContext(c.Request.Context(), &ctl.UserInfo{Id: claims.UserId}))

		c.Next()
	}
}

func extraceToken(ctx *gin.Context) string {
	bearerToken := ctx.GetHeader("Authorization")
	if bearerToken != "" {
		parts := strings.SplitN(bearerToken, " ", 2)

		if len(parts) == 2 || parts[0] == "Bearer" {

			return parts[1]
		}
		return bearerToken
	}

	return ""

}
