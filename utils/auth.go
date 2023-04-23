/**
 * Created by YuYoung on 2023/4/21
 * Description: auth 通用模块
 */

package utils

import (
	"github.com/dgrijalva/jwt-go"
	"short_link_sys_web_be/conf"
	"time"
)

func authInit() {
	Expire = time.Duration(conf.GlobalConfig.GetInt("auth.jwt.expire")) * time.Hour
	Secret = conf.GlobalConfig.GetString("auth.jwt.secret")
	Issuer = conf.GlobalConfig.GetString("auth.jwt.issuer")
}

var (
	Expire time.Duration
	Secret string
	Issuer string
)

type Token = jwt.StandardClaims

func GenerateToken() (string, error) {
	// token=头部+载荷+签名

	// 载荷=一些声明
	payload := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(Expire).Unix(),
		Issuer:    Issuer,
	}

	// 头部=算法+描述内容, 此处聚合头部和载荷
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	//签名并生成
	signedToken, err := token.SignedString([]byte(Secret))
	return signedToken, err
}

func ValidToken(tokenString string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Token{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	return token.Valid, err
}
