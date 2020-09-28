package jwt

import (
	"context"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

func TestJWT_verifyTokenSignature(t *testing.T) {
	type fields struct {
		staticSecrets        map[string]*config.Secret
		jwkSecrets           map[string]*jwkSecret
		mapJwkKidToSecretKid map[string]string
	}
	type args struct {
		ctx    context.Context
		token  string
		secret *config.Secret
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "HS256 JWT token with audience claim as a string value",
			fields: fields{
				staticSecrets: map[string]*config.Secret{"84095657997874536016": {
					IsPrimary: true,
					Alg:       config.HS256,
					KID:       "84095657997874536016",
					Audience:  []string{"myProject"},
					Secret:    "15557592078635013315",
				}},
			},
			args: args{
				ctx:   context.Background(),
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJteVByb2plY3QifQ.wVRCdD6QPF0qRWHWZt5z0vKWsnnia-tQrCxpsAnlwPk",
				secret: &config.Secret{
					IsPrimary: true,
					Alg:       config.HS256,
					KID:       "84095657997874536016",
					Audience:  []string{"myProject"},
					Secret:    "15557592078635013315",
				},
			},
			want:    map[string]interface{}{"aud": "myProject"},
			wantErr: false,
		},
		{
			name: "HS256 JWT token with audience claim as a array value",
			fields: fields{
				staticSecrets: map[string]*config.Secret{"84095657997874536016": {
					IsPrimary: true,
					Alg:       config.HS256,
					KID:       "84095657997874536016",
					Audience:  []string{"myProject2"},
					Secret:    "15557592078635013315",
				}},
			},
			args: args{
				ctx:   context.Background(),
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsibXlQcm9qZWN0MSIsIm15UHJvamVjdDIiXX0.cUmRkLzDWxZEq5etL8kpAU3JcFJVGfKQBxAw5j_CkpY",
				secret: &config.Secret{
					IsPrimary: true,
					Alg:       config.HS256,
					KID:       "84095657997874536016",
					Audience:  []string{"myProject2"},
					Secret:    "15557592078635013315",
				},
			},
			want:    map[string]interface{}{"aud": []interface{}{"myProject1", "myProject2"}},
			wantErr: false,
		},
		{
			name: "HS256 JWT token without audience claim but to validate token must contain audience",
			fields: fields{
				staticSecrets: map[string]*config.Secret{"84095657997874536016": {
					IsPrimary: true,
					Alg:       config.HS256,
					KID:       "84095657997874536016",
					Audience:  []string{"myProject"},
					Secret:    "15557592078635013315",
				}},
			},
			args: args{
				ctx:   context.Background(),
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.kGi2QNVvJ-w2JMI9N4rikFBYZSGyaZyD91APFsHT0_Y",
				secret: &config.Secret{
					IsPrimary: true,
					Alg:       config.HS256,
					KID:       "84095657997874536016",
					Audience:  []string{"myProject"},
					Secret:    "15557592078635013315",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "HS256 JWT token with empty audience claim but to validate token must contain audience",
			fields: fields{
				staticSecrets: map[string]*config.Secret{"84095657997874536016": {
					IsPrimary: true,
					Alg:       config.HS256,
					KID:       "84095657997874536016",
					Audience:  []string{"myProject"},
					Secret:    "15557592078635013315",
				}},
			},
			args: args{
				ctx:   context.Background(),
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIifQ.oSa-Z043KgT7v8JE504nhnxUb-uajVnBlMRKbuI7ipc",
				secret: &config.Secret{
					IsPrimary: true,
					Alg:       config.HS256,
					KID:       "84095657997874536016",
					Audience:  []string{"myProject"},
					Secret:    "15557592078635013315",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "HS256 JWT token with issuer claim as a string value",
			fields: fields{
				staticSecrets: map[string]*config.Secret{"84095657997874536016": {
					IsPrimary: true,
					Alg:       config.HS256,
					KID:       "84095657997874536016",
					Issuer:    []string{"https://my-auth-server.com"},
					Secret:    "15557592078635013315",
				}},
			},
			args: args{
				ctx:   context.Background(),
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL215LWF1dGgtc2VydmVyLmNvbSJ9.6Txu2dMelp1dc3b1D_KrlHcNvexswuDA1YHAB1Aculc",
				secret: &config.Secret{
					IsPrimary: true,
					Alg:       config.HS256,
					KID:       "84095657997874536016",
					Issuer:    []string{"https://my-auth-server.com"},
					Secret:    "15557592078635013315",
				},
			},
			want:    map[string]interface{}{"iss": "https://my-auth-server.com"},
			wantErr: false,
		},
		{
			name: "HS256 JWT token with issuer claim is an empty string value but required for validating token",
			fields: fields{
				staticSecrets: map[string]*config.Secret{"84095657997874536016": {
					IsPrimary: true,
					Alg:       config.HS256,
					KID:       "84095657997874536016",
					Issuer:    []string{"https://my-auth-server.com"},
					Secret:    "15557592078635013315",
				}},
			},
			args: args{
				ctx:   context.Background(),
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiIifQ.WkaZivZc325Z6NK3CF-DkQhTrZg09ly8y-EyQdRDtSw",
				secret: &config.Secret{
					IsPrimary: true,
					Alg:       config.HS256,
					KID:       "84095657997874536016",
					Issuer:    []string{"https://my-auth-server.com"},
					Secret:    "15557592078635013315",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JWT{
				staticSecrets:        tt.fields.staticSecrets,
				jwkSecrets:           tt.fields.jwkSecrets,
				mapJwkKidToSecretKid: tt.fields.mapJwkKidToSecretKid,
			}
			got, err := j.verifyTokenSignature(tt.args.ctx, tt.args.token, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("verifyTokenSignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("verifyTokenSignature() got = %v, want %v", got, tt.want)
			}
		})
	}
}
