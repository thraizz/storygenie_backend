package helper

import (
	"reflect"
	"testing"

	"github.com/golang-jwt/jwt/v4"
)

func TestVerifyJWT(t *testing.T) {
	type args struct {
		tokenString string
	}
	tests := []struct {
		name    string
		args    args
		want    jwt.MapClaims
		wantErr bool
	}{
		{
			name: "Valid token",
			args: args{
				tokenString: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o",
			},
			want: jwt.MapClaims{
				"sub":  "1234567890",
				"name": "John Doe",
				"iat":  float64(1516239022),
			},
			wantErr: false,
		},
		{
			name: "Invalid token",
			args: args{
				tokenString: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5d",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "No Bearer token",
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o",
			},
			want: jwt.MapClaims{
				"sub":  "1234567890",
				"name": "John Doe",
				"iat":  float64(1516239022),
			},
			wantErr: false,
		},
		{
			name: "Empty token",
			args: args{
				tokenString: "",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("JWT_PRIVATE_KEY", "secret")
			got, err := VerifyJWT(tt.args.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifyJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}
