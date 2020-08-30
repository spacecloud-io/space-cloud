// +build offline

package admin

import (
	"encoding/base64"
)

func init() {
	licenseMode = "offline"
	temp, _ := base64.RawStdEncoding.DecodeString(licensePublicKey)
	licensePublicKey = string(temp)
}
