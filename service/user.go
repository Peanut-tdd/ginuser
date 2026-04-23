package service

import (
	"context"
	"errors"
	"gin-user/config"
	"gin-user/pkg/jwt"
	"gin-user/pkg/utils/ctl"
	"gin-user/repository/db/model"
	"gin-user/types"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type userSrv struct {
}

var userSrvIn *userSrv

var syncOnce sync.Once

func NewUserSrv() *userSrv {

	syncOnce.Do(func() {
		userSrvIn = &userSrv{}
	})
	return userSrvIn
}

func (u *userSrv) RegisterUser(ctx context.Context, req *types.UserRegisterReq) (resp interface{}, err error) {

	userModel := model.NewUserModel(ctx)

	exist, err := userModel.GetUserByUsername(ctx, req.UserName)
	if err != nil {
		return nil, err
	}
	if len(exist) > 0 {
		err = errors.New("用户已存在")
		return
	}

	user := &model.User{
		NickName: req.NickName,
		Username: req.UserName,
		Password: req.Password,
		Age:      req.Age,
	}

	resp, err = userModel.Create(ctx, user)
	return resp, err

}

func (u *userSrv) LoginUser(ctx context.Context, req *types.UserLoginReq) (resp interface{}, err error) {
	ctx, span := otel.Tracer("gin-user/service/user").Start(ctx, "userSrv.LoginUser")
	defer span.End()
	span.SetAttributes(attribute.String("login.username", req.UserName))

	userModel := model.NewUserModel(ctx)

	user, err := userModel.GetUserByNamePasswd(ctx, req.UserName, req.Password)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "query user failed")
		return nil, err
	}

	if user != nil {
		userId := user.ID
		username := user.Username
		span.SetAttributes(attribute.Int64("user.id", int64(userId)))
		jm := jwt.NewJwtManager()

		tokenPair, err := jm.GenerateTokens(userId, username)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "generate token failed")
			return nil, err
		}
		span.SetStatus(codes.Ok, "login success")

		config.CtxInfof(ctx, "LoginUser resp:%v", resp)

		return tokenPair, nil

	}

	span.SetStatus(codes.Error, "invalid username or password")
	err = errors.New("用户名或密码错误")
	return

}

func (u *userSrv) UserRefresh(ctx context.Context) (resp interface{}, err error) {

	jm := jwt.NewJwtManager()
	user, err := ctl.GetUserInfo(ctx)

	if err != nil {
		return nil, err
	}

	refreshToken, err := jm.RefreshTokens(user.Id)
	if err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func (u *userSrv) UserLoginOut(authToken string) error {
	jm := jwt.NewJwtManager()

	err := jm.InValidateToken(authToken)
	if err != nil {
		return err
	}

	return nil
}
