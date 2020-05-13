package token

import (
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
					Msg:   "ok",
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

			gotResult := t.Encode(tt.args.appKey, tt.args.channel, tt.args.secret, tt.args.extra)

			t1.Logf("%+v", gotResult)
			//if gotResult := t.Encode(tt.args.appKey, tt.args.channel, tt.args.secret, tt.args.extra); !reflect.DeepEqual(gotResult, tt.wantResult) {
			//	t1.Errorf("Encode() = %v, want %v", gotResult, tt.wantResult)
			//}
		})
	}
}
