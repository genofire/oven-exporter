package helper

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"net/url"
)

type Policy struct {
	URLExpire    int    `json:"url_expire"`
	URLActivate  int    `json:"url_activate,omitempty"`
	StreamExpire int    `json:"stream_expire,omitempty"`
	AllowIP      string `json:"allow_ip,omitempty"`
}

func (p Policy) Encode() (string, error) {
	str, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(str), nil
}

func SignEncodedPolicy(u *url.URL, secretKey string) string {
	hasher := hmac.New(sha1.New, []byte(secretKey))
	hasher.Write([]byte(u.String()))
	return base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))

}

func (p Policy) SignWithQuery(u *url.URL, secretKey, encodeQuery string) (string, error) {
	encode, err := p.Encode()
	if err != nil {
		return "", nil
	}
	query := u.Query()
	query.Add(encodeQuery, encode)
	u.RawQuery = query.Encode()
	return SignEncodedPolicy(u, secretKey), nil
}
func (p Policy) Sign(u *url.URL, secretKey string) (string, error) {
	return p.SignWithQuery(u, secretKey, "policy")
}

func (p Policy) SignURLWithQuery(u *url.URL, secretKey, encodeQuery, signatureQuery string) error {
	encode, err := p.Encode()
	if err != nil {
		return err
	}
	query := u.Query()
	query.Add(encodeQuery, encode)
	u.RawQuery = query.Encode()

	signature := SignEncodedPolicy(u, secretKey)
	query.Add(signatureQuery, signature)
	u.RawQuery = query.Encode()
	return nil
}
func (p Policy) SignURL(u *url.URL, secretKey string) error {
	return p.SignURLWithQuery(u, secretKey, "policy", "signature")
}
