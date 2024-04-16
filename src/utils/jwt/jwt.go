package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"time"

	"github.com/dvsekhvalnov/jose2go/base64url"
)

type Header struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
}

type Payload struct {
	Exp   uint64 `json:"exp"`
	Email string `json:"email"`
}

type Token struct {
	Header  Header  `json:"header"`
	Payload Payload `json:"payload"`
	Sign    string  `json:"sign"`
}

func New(email, secret string) Token {
	now := time.Now().UnixMilli() + 60 * 60 * 1000
	t := Token{Header: Header{Typ: "JWT", Alg: "HS256"}, Payload: Payload{Exp: uint64(now), Email: email}, Sign: ""}

	t.createSign(secret)

	return t
}

func EncodeBase64Url(s []byte) string {
	return base64url.Encode(s)
}

func EncodeSha256(s []byte) []byte {
	hash := sha256.New()
	hash.Write(s)
	return hash.Sum(nil)
}

func EncodeHmac(s []byte, secret string) []byte {
	hm := hmac.New(sha256.New, []byte(secret))
	hm.Write(s)
	return hm.Sum(nil)
}

func (t *Token) createSign(secret string) {
	jsonHeader, _ := json.Marshal(t.Header)
	jsonPayload, _ := json.Marshal(t.Payload)
	dataString := EncodeBase64Url(jsonHeader) + "." + EncodeBase64Url(jsonPayload)
	t.Sign = EncodeBase64Url(EncodeHmac([]byte(dataString), secret))
}

func (t *Token) ToString() string {
	jsonHeader, _ := json.Marshal(t.Header)
	jsonPayload, _ := json.Marshal(t.Payload)
	dataString := EncodeBase64Url(jsonHeader) + "." + EncodeBase64Url(jsonPayload)

	return dataString + "." + t.Sign
}