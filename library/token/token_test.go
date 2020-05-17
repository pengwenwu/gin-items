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
			name: "success",
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
				secret:  "45f25874aa6dd33427dee744f2a800e6",
				extra: EncodeExtraData{
					LoginUserId: 1535917,
					NickName:    "四个二带俩王",
					BabyInfo: []map[string]interface{}{
						{
							"name":   "张三",
							"gender": "男",
							"age":    1,
						},
						{
							"name":   "李四",
							"gender": "女",
							"age":    1.5,
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
			name: "fail: invalid appKey",
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
				appKey:  "",
				channel: 0,
				secret:  "",
				extra: EncodeExtraData{
					LoginUserId: 1535917,
					NickName:    "四个二带俩王",
					BabyInfo: []map[string]interface{}{
						{
							"name":   "张三",
							"gender": "男",
							"age":    1,
						},
						{
							"name":   "李四",
							"gender": "女",
							"age":    1.5,
						},
					},
				},
			},
			wantResult: EncodeResult{
				Result: Result{
					State: 2001,
					Msg:   "",
				},
				Token: "",
			},
		},
		{
			name: "fail: invalid secret",
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
				extra: EncodeExtraData{
					LoginUserId: 1535917,
					NickName:    "四个二带俩王",
					BabyInfo: []map[string]interface{}{
						{
							"name":   "张三",
							"gender": "男",
							"age":    1,
						},
						{
							"name":   "李四",
							"gender": "女",
							"age":    1.5,
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
			if encodeResult.State != tt.wantResult.State {
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
		tokenString  string
		secret string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult DecodeResult
	}{
		{
			name: "success",
			fields: struct {
				expire int64
				before int64
			}{
				expire: 3600,
				before: 60,
			},
			args: struct {
				tokenString  string
				secret string
			}{
				tokenString:  "",
				secret: "45f25874aa6dd33427dee744f2a800e6",
			},
			wantResult: DecodeResult{
				Result: Result{
					State: 1,
					Msg:   "",
				},
				Data: &MyCustomClaims{
					AppKey:  "900ffe093ae07a09a99525baac3cfe53",
					Channel: 0,
					EncodeExtraData: EncodeExtraData{
						LoginUserId: 1535917,
						NickName:    "四个二带俩王",
						BabyInfo: []map[string]interface{}{
							{
								"name":   "张三",
								"gender": "男",
								"age":    1,
							},
							{
								"name":   "李四",
								"gender": "女",
								"age":    1.5,
							},
						},
					},
					StandardClaims: jwt.StandardClaims{},
				},
			},
		},
		{
			name: "fail: invalid tokenString",
			fields: struct {
				expire int64
				before int64
			}{
				expire: 3600,
				before: 1,
			},
			args: struct {
				tokenString  string
				secret string
			}{
				tokenString:  "",
				secret: "45f25874aa6dd33427dee744f2a800e6",
			},
			wantResult: DecodeResult{
				Result: Result{
					State: 2001,
					Msg:   "",
				},
				Data: nil,
			},
		},
		{
			name: "fail: invalid secret",
			fields: struct {
				expire int64
				before int64
			}{
				expire: 0,
				before: 0,
			},
			args: struct {
				tokenString  string
				secret string
			}{
				tokenString:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwidXNlcm5hbWUiOiJKb2huIERvZSIsImlhdCI6MTUxNjIzOTAyMiwiYXBwa2V5IjoiNWUyZjJjMzU2NTM4OTMyOWMxMjQ3ZWZkMDQzZjNiZTAuaW9zIn0.tOfZmMANQKth6oFVJqUT_LtMAxBmUr1BkFuhBmNS1E8",
				secret: "",
			},
			wantResult: DecodeResult{
				Result: Result{
					State: 2002,
					Msg:   "",
				},
				Data: nil,
			},
		},
		{
			name: "fail: token error",
			fields: struct {
				expire int64
				before int64
			}{
				expire: 0,
				before: 0,
			},
			args: struct {
				tokenString  string
				secret string
			}{
				tokenString:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBrZXkiOiI5MDBmZmUwOTNhZTA3YTA5YTk5NTI1YmFhYzNjZmU1MyIsImNoYW5uZWwiOjAsImxvZ2luX3VzZXJfaWQiOjE1MzU5MTcsIm5pY2tfbmFtZSI6IuWbm-S4quS6jOW4puS_qeeOiyIsImJhYnlfaW5mbyI6W3siYWdlIjoxLCJnZW5kZXIiOiLnlLciLCJuYW1lIjoi5byg5LiJIn0seyJhZ2UiOjEuNSwiZ2VuZGVyIjoi5aWzIiwibmFtZSI6IuadjuWbmyJ9XSwiYXVkIjoiY29tbW9uIiwiZXhwIjoxNTg5MTAwNzUyLCJpYXQiOjE1ODk3MDU1NTIsImlzcyI6ImFwaSIsIm5iZiI6MTU4OTcwNTU1Miwic3ViIjoiand.t87c6p6jMRR1XnAiYxCxizU-O0gUOVth8r0CzzR0cUo",
				secret: "45f25874aa6dd33427dee744f2a800e6",
			},
			wantResult: DecodeResult{
				Result: Result{
					State: 3001,
					Msg:   "",
				},
				Data: &MyCustomClaims{
					AppKey:  "900ffe093ae07a09a99525baac3cfe53",
					Channel: 0,
					EncodeExtraData: EncodeExtraData{
						LoginUserId: 1535917,
						NickName:    "四个二带俩王",
						BabyInfo: []map[string]interface{}{
							{
								"name":   "张三",
								"gender": "男",
								"age":    1,
							},
							{
								"name":   "李四",
								"gender": "女",
								"age":    1.5,
							},
						},
					},
					StandardClaims: jwt.StandardClaims{},
				},
			},
		},
		{
			name: "fail: token expired",
			fields: struct {
				expire int64
				before int64
			}{
				expire: 0,
				before: 0,
			},
			args: struct {
				tokenString  string
				secret string
			}{
				tokenString:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBrZXkiOiI5MDBmZmUwOTNhZTA3YTA5YTk5NTI1YmFhYzNjZmU1MyIsImNoYW5uZWwiOjAsImxvZ2luX3VzZXJfaWQiOjE1MzU5MTcsIm5pY2tfbmFtZSI6IuWbm-S4quS6jOW4puS_qeeOiyIsImJhYnlfaW5mbyI6W3siYWdlIjoxLCJnZW5kZXIiOiLnlLciLCJuYW1lIjoi5byg5LiJIn0seyJhZ2UiOjEuNSwiZ2VuZGVyIjoi5aWzIiwibmFtZSI6IuadjuWbmyJ9XSwiYXVkIjoiY29tbW9uIiwiZXhwIjoxNTg5MTAwNzUyLCJpYXQiOjE1ODk3MDU1NTIsImlzcyI6ImFwaSIsIm5iZiI6MTU4OTcwNTU1Miwic3ViIjoiand0In0.t87c6p6jMRR1XnAiYxCxizU-O0gUOVth8r0CzzR0cUo",
				secret: "45f25874aa6dd33427dee744f2a800e6",
			},
			wantResult: DecodeResult{
				Result: Result{
					State: 3002,
					Msg:   "",
				},
				Data: nil,
			},
		},

		{
			name: "fail: parse error",
			fields: struct {
				expire int64
				before int64
			}{
				expire: 0,
				before: 0,
			},
			args: struct {
				tokenString  string
				secret string
			}{
				tokenString:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwidXNlcm5hbWUiOiJKb2huIERvZSIsImlhdCI6MTUxNjIzOTAyMiwiYXBwa2V5IjoiNWUyZjJjMzU2NTM4OTMyOWMxMjQ3ZWZkMDQzZjNiZTAuaW9zIn0.tOfZmMANQKth6oFVJqUT_LtMAxBmUr1BkFuhBmNS1E8",
				secret: "45f25874aa6dd33427dee744f2a800e6",
			},
			wantResult: DecodeResult{
				Result: Result{
					State: 3003,
					Msg:   "",
				},
				Data: nil,
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			appKey := "900ffe093ae07a09a99525baac3cfe53"
			secret := "45f25874aa6dd33427dee744f2a800e6"

			var tokenString string
			t := NewToken()
			if tt.name == "success" {
				t.SetExpire(tt.fields.expire)
				t.SetBefore(tt.fields.before)
				encodeResult := t.Encode(appKey, 0, secret, tt.wantResult.Data.EncodeExtraData)
				tokenString = encodeResult.Token
			} else {
				tokenString = tt.args.tokenString
			}
			decodeResult := Decode(tokenString, tt.args.secret)

			if decodeResult.State != tt.wantResult.State {
				t1.Errorf("Decode() = %v, want %v", decodeResult.State, tt.wantResult.State)
				return
			}
			if decodeResult.Data == nil {
				return
			}
			if decodeResult.Data.AppKey != tt.wantResult.Data.AppKey {
				t1.Errorf("Decode() = %v, want %v", decodeResult.Data.AppKey, tt.wantResult.Data.AppKey)
				return
			}
			if decodeResult.Data.Channel != tt.wantResult.Data.Channel {
				t1.Errorf("Decode() = %v, want %v", decodeResult.Data.Channel, tt.wantResult.Data.Channel)
				return
			}
		})
	}
}

func Test_token_UnSafeDecode(t1 *testing.T) {
	type fields struct {
		expire int64
		before int64
	}
	type args struct {
		tokenString string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult DecodeResult
	}{
		{
			name:   "success",
			fields: fields{},
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwidXNlcm5hbWUiOiJKb2huIERvZSIsImlhdCI6MTUxNjIzOTAyMiwiYXBwa2V5IjoiNWUyZjJjMzU2NTM4OTMyOWMxMjQ3ZWZkMDQzZjNiZTAuaW9zIn0.tOfZmMANQKth6oFVJqUT_LtMAxBmUr1BkFuhBmNS1E8",
			},
			wantResult: DecodeResult{
				Result: Result{
					State: 1,
					Msg:   "",
				},
				Data: nil,
			},
		},
		{
			name:   "fail: error token",
			fields: fields{},
			args: args{
				tokenString: "",
			},
			wantResult: DecodeResult{
				Result: Result{
					State: 2001,
					Msg:   "",
				},
				Data: nil,
			},
		},
		{
			name:   "fail: error payload",
			fields: fields{},
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwidXNlcm5hbWUiOiJKb2huIERvZSIsImlhdCI6MTUxNjIzOTAyMiwiYXBwa2V5IjoiNWUyZjJjMzU2NTM4OTMyOWMxMjQ3ZWZkMDQzZjNiZTAuaW9zIn0=.tOfZmMANQKth6oFVJqUT_LtMAxBmUr1BkFuhBmNS1E8",
			},
			wantResult: DecodeResult{
				Result: Result{
					State: 2002,
					Msg:   "",
				},
				Data: nil,
			},
		},
		{
			name:   "fail: error data",
			fields: fields{},
			args: args{
				tokenString: "aa.bb.cc",
			},
			wantResult: DecodeResult{
				Result:         Result{
					State: 2003,
					Msg:   "",
				},
				Data: nil,
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			gotResult := UnSafeDecode(tt.args.tokenString)
			if gotResult.State != tt.wantResult.State {
				t1.Errorf("UnSafeDecode() = %v, want %v", gotResult.State, tt.wantResult.State)
				return
			}
		})
	}
}
