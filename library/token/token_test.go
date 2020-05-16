package token

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
)

func Test_token_Encode(t1 *testing.T) {
	type fields struct {
		expire int64
		before int64
	}
	type args struct {
		appKey  string
		channel int
		secret  string
		extra   EncodeExtraData
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult EncodeResult
	}{
		{
			name: "encode success",
			fields: struct {
				expire int64
				before int64
			}{
				expire: 3600,
				before: 3600,
			},
			args: struct {
				appKey  string
				channel int
				secret  string
				extra EncodeExtraData
			}{
				appKey:  "900ffe093ae07a09a99525baac3cfe53",
				channel: 0,
				secret:  "45f25874aa6dd33427dee744f2a800e6",
				extra:EncodeExtraData{
					LoginUserId: 1535917,
					NickName:    "四个二带俩王",
					BabyInfo:    []map[string]interface{}{
						{
							"name": "张三",
							"gender": "男",
							"age": 1,
						},
						{
							"name": "李四",
							"gender": "女",
							"age": 1.5,
						},
					},
				},
			},
			wantResult: EncodeResult{
				Result: Result{
					State: 1,
					Msg:   "",
				},
				Token: "",
			},
		},
		{
			name: "encode fail",
			fields: struct {
				expire int64
				before int64
			}{
				expire: 3600,
				before: 3600,
			},
			args: struct {
				appKey  string
				channel int
				secret  string
				extra   EncodeExtraData
			}{
				appKey:  "900ffe093ae07a09a99525baac3cfe53",
				channel: 0,
				secret:  "",
				extra:EncodeExtraData{
					LoginUserId: 1535917,
					NickName:    "四个二带俩王",
					BabyInfo:    []map[string]interface{}{
						{
							"name": "张三",
							"gender": "男",
							"age": 1,
						},
						{
							"name": "李四",
							"gender": "女",
							"age": 1.5,
						},
					},
				},
			},
			wantResult: EncodeResult{
				Result: Result{
					State: 2002,
					Msg:   "",
				},
				Token: "",
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := NewToken()
			t.SetExpire(tt.fields.expire)
			t.SetBefore(tt.fields.before)

			encodeResult := t.Encode(tt.args.appKey, tt.args.channel, tt.args.secret, tt.args.extra)
			if encodeResult.State != tt.wantResult.State{
				t1.Errorf("Encode err msg: %v", encodeResult.Msg)
			}
		})
	}
}

func Test_token_Decode(t1 *testing.T) {
	type fields struct {
		expire int64
		before int64
	}
	type args struct {
		token  string
		secret string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult DecodeResult
	}{
		{
			name: "decode success",
			fields: struct {
				expire int64
				before int64
			}{
				expire: 3600,
				before: 60,
			},
			args: struct {
				token  string
				secret string
			}{
				token: "",
				secret: "45f25874aa6dd33427dee744f2a800e6",
			},
			wantResult: DecodeResult{
				Result:         Result{
					State: 1,
					Msg:   "",
				},
				MyCustomClaims: &MyCustomClaims{
					AppKey:          "900ffe093ae07a09a99525baac3cfe53",
					Channel:         0,
					EncodeExtraData: EncodeExtraData{
						LoginUserId: 1535917,
						NickName:    "四个二带俩王",
						BabyInfo:    []map[string]interface{}{
							{
								"name": "张三",
								"gender": "男",
								"age": 1,
							},
							{
								"name": "李四",
								"gender": "女",
								"age": 1.5,
							},
						},
					},
					StandardClaims:  jwt.StandardClaims{},
				},
			},
		},
		{
			name: "decode fail for expired",
			fields: struct {
				expire int64
				before int64
			}{
				expire: -604800,
				before: 1,
			},
			args: struct {
				token  string
				secret string
			}{
				token: "",
				secret: "45f25874aa6dd33427dee744f2a800e6",
			},
			wantResult: DecodeResult{
				Result:         Result{
					State: 3002,
					Msg:   "",
				},
				MyCustomClaims: &MyCustomClaims{
					AppKey:          "900ffe093ae07a09a99525baac3cfe53",
					Channel:         0,
					EncodeExtraData: EncodeExtraData{
						LoginUserId: 1535917,
						NickName:    "四个二带俩王",
						BabyInfo:    []map[string]interface{}{
							{
								"name": "张三",
								"gender": "男",
								"age": 1,
							},
							{
								"name": "李四",
								"gender": "女",
								"age": 1.5,
							},
						},
					},
					StandardClaims:  jwt.StandardClaims{},
				},
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			appKey := "900ffe093ae07a09a99525baac3cfe53"
			secret := "45f25874aa6dd33427dee744f2a800e6"

			t := NewToken()
			t.SetExpire(tt.fields.expire)
			t.SetBefore(tt.fields.before)
			encodeResult := t.Encode(appKey, 0, secret, tt.wantResult.EncodeExtraData)
			decodeResult := t.Decode(encodeResult.Token, tt.args.secret)

			if decodeResult.State != tt.wantResult.State {
				t1.Errorf("Decode() = %v, want %v", decodeResult.State, tt.wantResult.State)
				return
			}
			if decodeResult.MyCustomClaims == nil {
				return
			}
			if decodeResult.AppKey != tt.wantResult.AppKey {
				t1.Errorf("Decode() = %v, want %v", decodeResult.AppKey, tt.wantResult.AppKey)
				return
			}
			if decodeResult.Channel != tt.wantResult.Channel {
				t1.Errorf("Decode() = %v, want %v", decodeResult.Channel, tt.wantResult.Channel)
				return
			}
		})
	}
}