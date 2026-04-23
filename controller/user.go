package controller

import (
	"errors"
	"fmt"
	"gin-user/config"
	"gin-user/pkg/e"
	"gin-user/pkg/utils/ctl"
	"gin-user/service"
	"gin-user/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func UserRegisterHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var req *types.UserRegisterReq
		if err := ctx.ShouldBind(&req); err != nil {
			ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, err.Error(), e.ERROR))
			return
		}

		if req.NickName == "" || req.UserName == "" || req.Password == "" || req.Age == 0 {
			err := errors.New("请求参数错误啦")
			ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, err.Error(), e.ERROR))
			return
		}

		fmt.Println(req.NickName, req.UserName, req.Password, req.Age)

		usr := service.NewUserSrv()

		resp, err := usr.RegisterUser(ctx.Request.Context(), req)
		if err != nil {
			ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, err.Error(), e.ERROR))
			return
		}

		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

func UserLoginUserHandler() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		var req *types.UserLoginReq

		if err := ctx.ShouldBind(&req); err != nil {
			ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, err.Error(), e.ERROR))
			return
		}

		usr := service.NewUserSrv()

		resp, err := usr.LoginUser(ctx.Request.Context(), req)

		if err != nil {
			ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, err.Error(), e.ERROR))
			return
		}

		config.CtxInfof(ctx.Request.Context(), "UserLoginUserHandler resp:%v", resp)

		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))

	}
}

func UserRefreshTokenHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		usr := service.NewUserSrv()

		resp, err := usr.UserRefresh(ctx.Request.Context())

		if err != nil {
			ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, err.Error(), e.ERROR))
			return
		}

		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))

	}
}

func UserLogoutHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authToken := ctx.GetHeader("Authorization")

		parts := strings.SplitN(authToken, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			err := errors.New("无效的token令牌")
			ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, err.Error(), e.ERROR))
			return
		}

		usr := service.NewUserSrv()
		err := usr.UserLoginOut(parts[1])
		if err != nil {
			ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, err.Error(), e.ERROR))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, nil))
	}
}

func UserInfoHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//usr := service.NewUserSrv()

	}
}
