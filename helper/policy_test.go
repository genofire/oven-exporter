// https://airensoft.gitbook.io/ovenmediaengine/access-control/signedpolicy
package helper

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	examplePolicyEncode              = "eyJ1cmxfZXhwaXJlIjoxMzk5NzIxNTgxfQ"
	exampleSecretKey                 = "1kU^b6"
	exampleURL                       = "ws://192.168.0.100:3333/app/stream"
	exampleSignature                 = "dvVdBpoxAeCPl94Kt5RoiqLI0YE"
	exampleURLWithSignatureAndPolicy = "ws://192.168.0.100/app/stream?policy=eyJ1cmxfZXhwaXJlIjoxMzk5NzIxNTgxfQ&signature=dvVdBpoxAeCPl94Kt5RoiqLI0YE"
)

var (
	examplePolicy = Policy{
		URLExpire: 1399721581,
	}
)

func TestPolicyEncode(t *testing.T) {
	assert := assert.New(t)

	encode, err := examplePolicy.Encode()
	assert.NoError(err)
	assert.Equal(examplePolicyEncode, encode)
}

func TestPolicySign(t *testing.T) {
	assert := assert.New(t)

	u, err := url.Parse(exampleURL)
	assert.NoError(err)

	sign, err := examplePolicy.Sign(u, exampleSecretKey)
	assert.NoError(err)
	assert.Equal(exampleSignature, sign)
}

func TestPolicySignURL(t *testing.T) {
	assert := assert.New(t)

	u, err := url.Parse(exampleURL)
	assert.NoError(err)

	err = examplePolicy.SignURL(u, exampleSecretKey)
	assert.NoError(err)

	// drop port -> is not part of example
	u.Host = u.Hostname()

	assert.Equal(exampleURLWithSignatureAndPolicy, u.String())
}
