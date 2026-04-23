package jwt

import (
	"context"
	"errors"
	"fmt"
	"gin-user/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"sync"
	"time"
)

type JwtManager struct {
	Redis         *redis.Client
	Secret        string        `json:"secret"`
	AccessExpire  time.Duration `json:"access_expire"`
	RefreshExpire time.Duration `json:"refresh_expire"`
}

type UserClaims struct {
	UserId   uint   `json:"user_id"`
	Username string `json:"user_name"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

var jwtManager *JwtManager
var synconce sync.Once

func NewJwtManager() *JwtManager {
	synconce.Do(func() {
		jwtManager = &JwtManager{
			Redis:         config.RedisClient,
			Secret:        config.Conf.Jwt.Secret,
			AccessExpire:  time.Duration(config.Conf.Jwt.AccessExpire) * time.Minute,
			RefreshExpire: time.Duration(config.Conf.Jwt.RefreshExpire) * time.Hour,
		}
	})
	return jwtManager
}

func (jm *JwtManager) GenerateTokens(userId uint, username string) (*TokenPair, error) {
	accessToken, err := jm.generateAccessToken(userId, username)

	if err != nil {
		return nil, err
	}

	refreshToken, err := jm.RefreshTokens(userId)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    cast.ToInt(jm.AccessExpire.Seconds()),
	}, nil

}

func (jm *JwtManager) generateAccessToken(userId uint, username string) (string, error) {

	claims := &UserClaims{
		UserId:   userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jm.AccessExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.Secret))

}

func (jm *JwtManager) RefreshTokens(userId uint) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(jm.RefreshExpire)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        cast.ToString(userId),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.Secret))
}

func (jm *JwtManager) ParseToken(tokenString string) (*UserClaims, error) {

	blacklisted, err := jm.Redis.Exists(context.Background(), fmt.Sprintf("blacklist:%v", tokenString)).Result()

	if err != nil {
		return nil, err
	}
	if blacklisted > 0 {
		return nil, errors.New("token is blacklisted")
	}

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return []byte(jm.Secret), nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {

		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}

func (jm *JwtManager) InValidateToken(tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return []byte(jm.Secret), nil
	})

	if err != nil {
		return err
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {

		expireAt := claims.RegisteredClaims.ExpiresAt.Time
		ttl := time.Until(expireAt)

		if ttl <= 0 {
			return nil
		}

		err := jm.Redis.Set(context.Background(), fmt.Sprintf("blacklist:%v", tokenString), "1", ttl).Err()

		return err
	}

	return errors.New("无效token")

}
