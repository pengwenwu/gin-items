package token

import (
	"reflect"
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
		extra   map[string]interface{}
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
				extra   map[string]interface{}
			}{
				appKey:  "900ffe093ae07a09a99525baac3cfe53",
				channel: 0,
				secret:  "45f25874aa6dd33427dee744f2a800e6",
				extra: map[string]interface{}{
					"login_user_id": 1535917,
					"nickname":      "暗夜御林",
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
				extra   map[string]interface{}
			}{
				appKey:  "900ffe093ae07a09a99525baac3cfe53",
				channel: 0,
				secret:  "",
				extra: map[string]interface{}{
					"login_user_id": 1535917,
					"nickname":      "暗夜御林",
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
			},
			args: struct {
				token  string
				secret string
			}{
			},
			wantResult: DecodeResult{
				Result: Result{
					State: 1,
					Msg:   "ok",
				},
				Data: map[string]interface{}{
					"login_user_id": 1535917,
					"nickname": "四个二带俩王",
					"baby_info": []map[string]interface{}{
						{
							"realname": "张三",
							"gender": "男",
							"age": 1,
						},
						{
							"realname": "李四",
							"gender": "女",
							"age": 2,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			secret := "45f25874aa6dd33427dee744f2a800e6"
			appKey := "900ffe093ae07a09a99525baac3cfe53"

			t := NewToken()
			t.SetExpire(3600)
			t.SetBefore(3600)
			encodeResult := t.Encode(appKey, 0, secret, tt.wantResult.Data)
			decodeResult := t.Decode(encodeResult.Token, secret)
			t1.Logf("%+v %+v", encodeResult, decodeResult)

			if decodeResult.State != tt.wantResult.State ||
				!reflect.DeepEqual(decodeResult.Data, tt.wantResult.Data) {
				t1.Errorf("Decode() = %v, want %v", decodeResult, tt.wantResult)
			}
			for k,v := range tt.wantResult.Data {
				decodeVal, ok := decodeResult.Data[k]
				if !ok || decodeVal != v {
					t1.Errorf("Decode data = %v, want data %v", decodeVal, v)
				}
			}
		})
	}
}