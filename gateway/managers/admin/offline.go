// +build offline

package admin

import (
	"encoding/base64"
)

func init() {
	licenseMode = licenseModeOffline
	temp, _ := base64.RawStdEncoding.DecodeString(licensePublicKey)
	licensePublicKey = string(temp)
}
