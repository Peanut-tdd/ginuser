package model

import (
	"context"
	"fmt"
	"gin-user/config"
	"gorm.io/gorm"
)

const UserShardCount = 3 // 假设有3个分表：users_0, users_1, users_2

type User struct {
	gorm.Model
	NickName string `gorm:"column:nickname;type:varchar(255);not null"`
	Username string `gorm:"unique;column:username;type:varchar(255);not null"`
	Password string `gorm:"column:password;not null"`
	Age      int    `gorm:"column:age;not null"`
}

type UserModel struct {
	*gorm.DB
}

func NewUserModel(ctx context.Context) *UserModel {
	return &UserModel{config.NewDBClient(ctx)}
}

// GetTableName 根据分表键（如用户ID）获取表名
func GetTableName(userId uint) string {
	return fmt.Sprintf("users_%d", userId%UserShardCount)
}

func (u *UserModel) Create(ctx context.Context, user *User) (*User, error) {
	// 在插入前，根据某个分表键确定表名
	// 注意：如果是自增ID，插入前可能没有ID，通常分表键会选择业务上的唯一标识（如 Username 的哈希值）
	//tableName := GetTableName(user.ID)
	//result := u.DB.Table(tableName).Create(user)
	result := u.DB.WithContext(ctx).Create(user)

	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// GetUserByID 查询分表示例
func (u *UserModel) GetUserByID(ctx context.Context, userId uint) (*User, error) {
	var user User
	//tableName := GetTableName(userId)
	//err := u.DB.Table(tableName).Where("id = ?", userId).First(&user).Error

	err := u.DB.WithContext(ctx).Where("id=?", userId).First(&user).Error

	return &user, err
}

func (u *UserModel) GetUserByUsername(ctx context.Context, username string) ([]User, error) {
	var users []User
	result := u.DB.WithContext(ctx).Where("username=?", username).Find(&users)
	//b, _ := json.Marshal(users)
	//fmt.Printf("GetUserByUsername,user:%v, error : %s\n", string(b), result.Error)
	return users, result.Error
}

func (u *UserModel) GetUserByNamePasswd(ctx context.Context, username, passwd string) (*User, error) {
	var user User

	result := u.DB.WithContext(ctx).Where("username=? and password=?", username, passwd).Find(&user)

	return &user, result.Error

}
