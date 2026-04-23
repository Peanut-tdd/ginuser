package ctl

import (
	"context"
	"errors"
)

var UserKey string

type UserInfo struct {
	Id uint `json:"id"`
}

func GetUserInfo(ctx context.Context) (*UserInfo, error) {
	user, ok := FromContext(ctx)
	if !ok {
		return nil, errors.New("获取用户信息失败")
	}
	return user, nil

}

func NewContext(ctx context.Context, user *UserInfo) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

func FromContext(ctx context.Context) (*UserInfo, bool) {
	u, ok := ctx.Value(UserKey).(*UserInfo)
	return u, ok
}
