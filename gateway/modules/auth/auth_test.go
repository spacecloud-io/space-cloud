package auth

import (
	"context"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
)

// this test generates a tokenstring even if object is empty is this the behaviour we want
func TestCreateToken(t *testing.T) {
	var authCreateToken = []struct {
		testName       string
		secretKeys     []*config.Secret
		IsTokenInvalid bool
		IsErrExpected  bool
		object         map[string]interface{}
	}{
		{testName: "Successful Test", IsTokenInvalid: false, secretKeys: []*config.Secret{{IsPrimary: true, Secret: "mySecretkey"}}, object: map[string]interface{}{"id": "internal-sc-user"}},
		{testName: "Test Case-Invalid Token", IsTokenInvalid: true, IsErrExpected: false, secretKeys: []*config.Secret{{IsPrimary: true, Secret: "mySecretkey"}}, object: map[string]interface{}{"id": "internal-scuser"}},
		{testName: "Invalid Test Case-Empty Object", IsTokenInvalid: true, IsErrExpected: false, secretKeys: []*config.Secret{{IsPrimary: true, Secret: "mySecretkey"}}},
	}
	authModule := Init("chicago", "1", &crud.Module{}, nil)
	for _, test := range authCreateToken {
		t.Run(test.testName, func(t *testing.T) {
			_ = authModule.SetConfig(context.TODO(), "local", &config.ProjectConfig{Secrets: test.secretKeys}, nil, nil, nil, nil, config.EventingRules{}, config.SecurityFunctions{})
			_, err := authModule.CreateToken(context.Background(), test.object)
			if (err != nil) != test.IsErrExpected {
				t.Error("Got Error", err, "Wanted Error ", test.IsErrExpected)
			}
		})
	}
}

func TestIsTokenInternal(t *testing.T) {
	var authCreateToken = []struct {
		testName, token string
		secretKeys      []*config.Secret
		IsErrExpected   bool
	}{
		{testName: "Unsuccessful Test-Token has not been internally created", secretKeys: []*config.Secret{{IsPrimary: true, Secret: "mySecretkey"}}, IsErrExpected: true, token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc"},
		{testName: "Unsuccessful Test-Signature is Invalid", IsErrExpected: true, secretKeys: []*config.Secret{{IsPrimary: true, Secret: "mySecretkey"}}, token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.MKIZkrXy6nUMu5ejqiYKl7EOU1TxEoKTOww-eoQm6Lw"},
		{testName: "Successful Test Case", IsErrExpected: false, secretKeys: []*config.Secret{{IsPrimary: true, Secret: "mySecretkey"}}, token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImludGVybmFsLXNjLXVzZXIifQ.k3OcidcCnshBOGtzpprfV5Fhl2xWb6sjzPZH3omDDpw"},
	}

	authModule := Init("chicago", "1", &crud.Module{}, nil)
	for _, test := range authCreateToken {
		t.Run(test.testName, func(t *testing.T) {
			_ = authModule.SetConfig(context.TODO(), "local", &config.ProjectConfig{Secrets: test.secretKeys}, nil, nil, nil, nil, config.EventingRules{}, config.SecurityFunctions{})
			err := authModule.IsTokenInternal(context.Background(), test.token)
			if (err != nil) != test.IsErrExpected {
				t.Error("Got This ", err, "Wanted Error-", test.IsErrExpected)
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	var testCases = []struct {
		name          string
		IsErrExpected bool
		secretKeys    []*config.Secret
		token         string
		wantThis      map[string]interface{}
		reason        error
	}{
		{
			name:          "Single secret with Alg not set in config should parse token as HMAC token",
			IsErrExpected: false,
			secretKeys:    []*config.Secret{{IsPrimary: true, Secret: "50Au5e9DlwWzgTCeG7dWupvtb29TBJhkOd"}},
			token:         "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJ2ZXJzaW9uIjoidjEifQ.PLJ7bw6MTPuLflW--6EXT7-aVrkg5O9c9Af11D-p2N8",
			wantThis:      TokenClaims{"role": "admin", "version": "v1"},
		},
		{
			name:          "Single secret with HMAC256 Alg should parse token",
			IsErrExpected: false,
			secretKeys:    []*config.Secret{{IsPrimary: true, Alg: config.HS256, Secret: "50Au5e9DlwWzgTCeG7dWupvtb29TBJhkOd"}},
			token:         "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJ2ZXJzaW9uIjoidjEifQ.PLJ7bw6MTPuLflW--6EXT7-aVrkg5O9c9Af11D-p2N8",
			wantThis:      TokenClaims{"role": "admin", "version": "v1"},
		},
		{
			name:          "Multiple secret with HMAC256 Alg should parse token",
			IsErrExpected: false,
			secretKeys:    []*config.Secret{{IsPrimary: true, Alg: config.HS256, Secret: "50Au5e9DlwWzgTCeG7dWupvtb29TBJhkOd"}, {IsPrimary: false, Alg: config.HS256, Secret: "GEMKOjDiGQgqJXtiMOy3CG1V3EcVp0t8LzKL38THBBWPB2m"}},
			token:         "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJ2ZXJzaW9uIjoidjEifQ.wTcKB2Gp1xjxrn3XHhZ-angilidBpNXkjDKDKu2jckw",
			wantThis:      TokenClaims{"role": "admin", "version": "v1"},
		},
		{
			name:          "Single secret with RS256 Alg should parse token",
			IsErrExpected: false,
			secretKeys: []*config.Secret{
				{
					IsPrimary: true,
					Alg:       config.RS256,
					PrivateKey: `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA+GEMKOjDiGQgqJXtiMOy3CG1V3EcVp0t8LzKL38THBBWPB2m
Ghp2SMw2VjsLhgdmUeEhZ+9Btn8OvcE8jCpXUazbYiRNEoyWFiXOUPOI+2a8Ka4F
qmZ/Dv3Mq0D0Ybj+Db2QxdzYwdVXiZNSmkZVnwEzgRHs6oGCtEZBe11ZQgJx9W55
GUSU796j1LSrOLO0eadr0krFFtNtJLON70gZOUCVVGtAACTxd9T8M6DnioDIWeT6
E2ZllOCgRbGFUEjlfxpvXcJko1R13+u/2IdSASagEwagcvqKhYXbDyc4s2iechFq
Kq50Au5e9DlwWzgTCeG7dWupvtb29TBJhkOd/wIDAQABAoIBAQDX2qdRqi+8CaBk
Qp/DGqgHLPhG/NLu5vx02e0wZ9Q4sG3xIqcmylZ/n5GdUQULvsSKZge9Jq3KJxOt
jxDKk0V9gqbJ3NhqRWlt5/8sVOl36mmeI1FjLF5Brgm/ztxTTQzk7hiyC+MIWRsR
QcoorvVcERQsmGJ3zoPDncbsqkyW0ojhwxiStHwAlybdzxSXea0tt81Sxo5Yzinp
yO9jSpy11Hdq1VZxDt6nmdOUaEG17w5F/sFWJGvaGSwrbFispudfnErMpd+sE//A
x7llL2aWZrZ4/LXqOs1KVZ28Syg5UNGIZEhp5PTrS2AI2/1UjwbYdhJCrtVD9WX0
8W4MEGN5AoGBAP2DXDlBQVdb7vyr2KsedG5xrziRfo4Bkk348XDMpNgYH6p9Mvx6
lDUqZ76yGTs8XJjFQy4SbNnt3gBRjl4Xm2hlS38RxLp4rqUr1ilxXQtDexPaDClF
XjFDAqmJXYjEPByCmBNmcpwg+3CiYDG36qRWeBslL2rehCpBqUNra1y1AoGBAPrQ
y1fAnfl7Tv5+SGh6TKEMb/s3+HPD7+S3cpZCcDGZFHPEqpwRKtA0BSPyl6IsVASy
GYJgFdlh+ZEPTlDsh3Hcjz3IVfoFYOXI8mpDwr18G6JyM1fHJipIx1MwWymw7zsf
sLYBbIpgx/7PXFMy9DRtqJX3sTs8hruOgbhVgjRjAoGBAOt3v4lOmypRMcFuvGj0
eKC850tbmHi97O6Pc3KaIeK1RXRNpbXtBQv8vy/YrKbgggFD+AorAv/IeJNDLUEo
HCjzLRRxpJCu10iDVwBXU1zK6N6fgPw0ejC3VRmzhTUrT8oLYsViOsHZwLH1n/OO
KFtMFXLhLxHMbGXzZwxAEhChAoGBAJ7sacwCmpaYEWekMNEynRd2wAXYYy9HOdwc
eRjIpDppGtJ4DPqyzgP60j8C2h3CJMKM7yjzJSUGtZG6tw8DsJbvADxPklrHyawP
9bprkRtrZj86SVoXMBGe593ISBtUp2E5JUlOAa50wISuc3usT5xg12+e8MfuBBkX
pQ5d8BCtAoGADpEUTjkvOyMd8LRD09ASy0gHZk3XowtG1k9zkhh3wVbuaTOVurIX
crV25vkoIRmukN6rhCGLotiCQX0b9GihkBdEytJkIH0tei4+my1CywzK9tvBaUU/
aSV4YpIhfCydLcTmzb2Amd/99EfTC0gWwmlel2qRCltrXECB1bz5qMI=
-----END RSA PRIVATE KEY-----`,
					PublicKey: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+GEMKOjDiGQgqJXtiMOy
3CG1V3EcVp0t8LzKL38THBBWPB2mGhp2SMw2VjsLhgdmUeEhZ+9Btn8OvcE8jCpX
UazbYiRNEoyWFiXOUPOI+2a8Ka4FqmZ/Dv3Mq0D0Ybj+Db2QxdzYwdVXiZNSmkZV
nwEzgRHs6oGCtEZBe11ZQgJx9W55GUSU796j1LSrOLO0eadr0krFFtNtJLON70gZ
OUCVVGtAACTxd9T8M6DnioDIWeT6E2ZllOCgRbGFUEjlfxpvXcJko1R13+u/2IdS
ASagEwagcvqKhYXbDyc4s2iechFqKq50Au5e9DlwWzgTCeG7dWupvtb29TBJhkOd
/wIDAQAB
-----END PUBLIC KEY-----`,
				}},
			token:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJ2ZXJzaW9uIjoidjEifQ.0S2gD9pwbWL7e9pHX_vLSmplAu_N8eMtyvliEJ5KlkxJRd_YCKJ9ubM7Zb6mJHIZHwAOrz88gKnVpKQpWy9y4WRUvozd0XInt4FIXwPrE1xyRKyBWkC_7vFVMcbzSC_DQKaxh2N15ro9DaErcG8-HCDJ94IKydeEe9dVdMOWS45qQ5Kf2Jeyh7HDUvtADEspjAkrdp7lIGobd_9HVNQyrgTIEmBQtItgIxdX617oOFW4evSNm5OpMO6oUJ1OQtTdz-ci2lf3gRnMXCctSIIUFIzJZsOigicMUMucTTWAT_5ZWdIt2yeIUYavP2WUyc5MdZaK69Z0_qz40flJS5UKJw",
			wantThis: TokenClaims{"role": "admin", "version": "v1"},
		},
		{
			name:          "Test should fail if there is a mis match in signing method",
			IsErrExpected: true,
			secretKeys:    []*config.Secret{{IsPrimary: true, Secret: "mySecretkey"}},
			token:         "eyJhbGciOiJQUzM4NCIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.nakZ1JcYWHcXcG1ZfIY7mJNwcVPQ7U1HvuLEsG9fyz-H9ig3ql8BiI3T-7A2PHe-lBIxjS7hXx8O8lxMg7y7rqUHtPLAGOuCd4Ft88KupgPcF5w-KVpeSgWl598zNLWqJpjcwiPewt3gsU6pwSaTz24JmfZQRrDX8KOtejaGs5OECdk2dDW2rwO98npNX39yYx6eSfZbXCLJ7wIhT3UDbuaOGHnD3wyEtih013NDrnkvVXJRKXUwF7F-g31NWgEgVt-tWkR5vcBBSRYKzIbD7-wxpV4ifLp_XdbVNl3Uf7ja6FeUnGq1Pb9AnAY7lD4Rk7sYQe4P-ATHtkgSg5levw",
		},
		{
			name:          "Test should fail for an invalid token",
			IsErrExpected: true,
			secretKeys:    []*config.Secret{{IsPrimary: true, Secret: "mysecretkey"}},
			token:         "1234.abcd",
		},
		{
			name:          "Test should fail for invalid signature-illegal base64 data at input",
			IsErrExpected: true,
			secretKeys:    []*config.Secret{{IsPrimary: true, Secret: "mySecretkey"}}, token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N",
		},
		{
			name:          "Test should fail as invalid secret key-invalid signature",
			IsErrExpected: true,
			secretKeys:    []*config.Secret{{IsPrimary: true, Secret: "mysecretkey"}},
			token:         "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
		},
	}
	authModule := Init("chicago", "1", &crud.Module{}, nil)
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			if err := authModule.SetConfig(context.TODO(), "local", &config.ProjectConfig{Secrets: test.secretKeys}, nil, nil, nil, nil, nil, config.SecurityFunctions{}); err != nil {
				t.Errorf("error setting config of auth module  - %s", err.Error())
			}
			tokenClaims, err := authModule.jwt.ParseToken(context.Background(), test.token)
			if (err != nil) != test.IsErrExpected {
				t.Error(test.name, ": Got:", err, "Wanted Error:", test.IsErrExpected)
			}
			if !reflect.DeepEqual(test.wantThis, tokenClaims) {
				t.Error(test.name, ": Got:", tokenClaims, "Want:", test.wantThis, "Reason:", err)
			}

		})
	}
}
