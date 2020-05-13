package token

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const alg  = "HS256" // 签名算法

type token struct {
	expire int64 // 生成后，多少秒内有效
	before int64 // 生成前，多少秒有效（防止时间不同步）
}

type Result struct {
	State int `json:"state"`
	Msg string `json:"msg"`
}

type EncodeResult struct {
	Result
	Token string `json:"token"`
}

type DecodeResult struct {
	Result
	Data interface{} `json:"data"`
}

func NewToken() *token {
	return &token{
		expire: 604800,
		before:3600,
	}
}

func (t *token) SetExpire(expire int64) {
	t.expire = expire
}

func (t *token) SetBefore(before int64) {
	t.before = before
}

func (t *token) Encode(appKey string, channel int, secret string, extra map[string]interface{}) (result EncodeResult) {
	if len(appKey) < 32 {
		result.State = 2001
		result.Msg = "appkey是32位字符串，请传入正确的值"
		return
	}
	if len(secret) == 0 {
		result.State = 2002
		result.Msg = "请输入正确的secret"
		return
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["appkey"] = appKey
	claims["channel"] = channel
	claims["exp"] = time.Now().Unix() + t.expire
	claims["nbf"] = time.Now().Unix() - t.before
	for k, v := range extra {
		claims[k] = v
	}
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		result.State = 3001
		result.Msg = "生成token失败"
		return
	}
	result.State = 1
	result.Msg = "生成成功"
	result.Token = tokenString

	return result
}

func (t *token) Decode(token ,secret string) (result DecodeResult) {
	if len(token) == 0 {
		result.State = 2001
		result.Msg = "token错误，非空字符串"
		return
	}
	if len(secret) == 0 {
		result.State = 2002
		result.Msg = "secret错误，非空字符串"
		return
	}
	data, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if data.Valid {
		result.State = 1
		result.Msg = "解码成功"
		result.Data = data
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			result.State = 3001
			result.Msg = "token非法"
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			result.State = 3002
			result.Msg = "token尚未生效或者已过期"
		} else {
			result.State = 3003
			result.Msg = "token解析失败"
		}
	} else {
		result.State = 3003
		result.Msg = "token解析失败"
	}
	return

	return
}
