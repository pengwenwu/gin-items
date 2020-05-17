package token

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

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

type MyCustomClaims struct {
	AppKey string `json:"appkey"`
	Channel int `json:"channel"`
	EncodeExtraData
	jwt.StandardClaims
}

type EncodeExtraData struct {
	LoginUserId int `json:"login_user_id"`
	NickName string `json:"nick_name"`
	BabyInfo []map[string]interface{} `json:"baby_info"`
}

type DecodeResult struct {
	Result
	Data *MyCustomClaims `json:"data"`
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

func (t *token) Encode(appKey string, channel int, secret string, extra EncodeExtraData) (result EncodeResult) {
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
	now := time.Now().Unix()
	claims := MyCustomClaims{
		AppKey:          appKey,
		Channel:         channel,
		EncodeExtraData: extra,
		StandardClaims:  jwt.StandardClaims{
			IssuedAt:  now,
			ExpiresAt: now + t.expire,
			NotBefore: now - t.before,
			Issuer:    "api",
			Subject:   "jwt",
			Audience:  "common",
			Id:        "",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
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

func (t *token) Decode(tokenString ,secret string) (result DecodeResult) {
	if len(tokenString) == 0 {
		result.State = 2001
		result.Msg = "token错误，非空字符串"
		return
	}
	if len(secret) == 0 {
		result.State = 2002
		result.Msg = "secret错误，非空字符串"
		return
	}
	token , err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if token.Valid {
		claims, ok := token.Claims.(*MyCustomClaims)
		if !ok {
			result.State = 3003
			result.Msg = "token解析失败"
		} else {

		}
		result.State = 1
		result.Msg = "解码成功"
		result.Data = claims
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
}

// 只有在网关后，获取得到已验证过的token，将其值返回回来，不允许在网关前使用
func (t *token) UnSafeDecode(tokenString string) (result DecodeResult) {
	splitArr := strings.Split(tokenString, ".")
	if len(splitArr) != 3 {
		result.State = 2001
		result.Msg = "token由3个.号分割，格式不对"
		return
	}
	payloadStr, err := base64.RawStdEncoding.DecodeString(splitArr[1])
	if err != nil {
		result.State = 2002
		result.Msg = "token payload解析失败"
		return
	}
	if err := json.Unmarshal(payloadStr, &result.Data); err != nil {
		result.State = 2003
		result.Msg = "token不符合jwt规范，非法token"
		return
	}
	result.State = 1
	result.Msg = "解压成功，数据未经安全校验，请谨慎使用"
	return
}
