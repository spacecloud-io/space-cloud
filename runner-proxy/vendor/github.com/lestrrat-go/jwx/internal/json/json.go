package json

import (
	"bytes"
	"sync"

	"github.com/lestrrat-go/jwx/internal/base64"
	"github.com/pkg/errors"
)

var muGlobalConfig sync.RWMutex
var useNumber bool

// Sets the global configuration for json decoding
func DecoderSettings(inUseNumber bool) {
	muGlobalConfig.Lock()
	useNumber = inUseNumber
	muGlobalConfig.Unlock()
}

// Unmarshal respects the values specified in DecoderSettings,
// and uses a Decoder that has certain features turned on/off
func Unmarshal(b []byte, v interface{}) error {
	dec := NewDecoder(bytes.NewReader(b))
	return dec.Decode(v)
}

func AssignNextBytesToken(dst *[]byte, dec *Decoder) error {
	var val string
	if err := dec.Decode(&val); err != nil {
		return errors.Wrap(err, `error reading next value`)
	}

	buf, err := base64.DecodeString(val)
	if err != nil {
		return errors.Errorf(`expected base64 encoded []byte (%T)`, val)
	}
	*dst = buf
	return nil
}

func ReadNextStringToken(dec *Decoder) (string, error) {
	var val string
	if err := dec.Decode(&val); err != nil {
		return "", errors.Wrap(err, `error reading next value`)
	}
	return val, nil
}

func AssignNextStringToken(dst **string, dec *Decoder) error {
	val, err := ReadNextStringToken(dec)
	if err != nil {
		return err
	}
	*dst = &val
	return nil
}
