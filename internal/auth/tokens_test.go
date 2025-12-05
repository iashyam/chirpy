package auth_test

import (
	"chirpy/internal/auth"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {

	realUUD := uuid.New()
	secret := "Vedha is *****"
	toke_string, _ := auth.MakeJWT(realUUD, secret, time.Duration(time.Second*3))

	tests := []struct {
		name        string // description of this test case
		tokenString string
		tokenSecret string
		want        uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Everything Okay",
			tokenString: toke_string,
			tokenSecret: "Vedha is *****",
			want:        realUUD,
			wantErr:     false,
		},
		{
			name:        "Not Matching secret string",
			tokenString: toke_string,
			tokenSecret: "Vedha is a good person",
			want:        realUUD,
			wantErr:     true,
		},
		{
			name:        "Not Matching Token",
			tokenString: "Vedha is a good person",
			tokenSecret: "Vedha is *****",
			want:        realUUD,
			wantErr:     true,
		},
		{
			name:        "Expired Token",
			tokenString: toke_string,
			tokenSecret: "Vedha is *****",
			want:        realUUD,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Expired Token" {
				time.Sleep(time.Second * 4)
			}
			got, gotErr := auth.ValidateJWT(tt.tokenString, tt.tokenSecret)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ValidateJWT() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ValidateJWT() succeeded unexpectedly")
			}
			if tt.want.String() != got.String() {
				t.Errorf("ValidateJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}
