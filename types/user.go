package types

type UserRegisterReq struct {
	NickName string `form:"nick_name" json:"nick_name"`
	UserName string `form:"user_name" json:"user_name"`
	Password string `form:"password" json:"password"`
	Age      int    `form:"age" json:"age"`
}

type UserLoginReq struct {
	UserName string `form:"user_name" json:"user_name"`
	Password string `form:"password" json:"password"`
}
