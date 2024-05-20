package util

import (
	"awesomeProject/model"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

// RsaEncrypt encrypts data using rsa public key.
func RsaEncrypt(key string, data []byte) ([]byte, error) {
	strPEM := `
-----BEGIN PUBLIC KEY-----
` + key + `
-----END  PUBLIC KEY-----
`

	block, _ := pem.Decode([]byte(strPEM))
	if block == nil {
		return nil, errors.New("decode public key error")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.EncryptPKCS1v15(rand.Reader, pub.(*rsa.PublicKey), data)
}

// RsaDecrypt decrypts data using rsa private key.
func RsaDecrypt(ciphertext string) (string, error) {
	key := model.PrivateKey
	decodedtext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed, error=%s\n", err.Error())
	}
	v15, err := rsa.DecryptPKCS1v15(rand.Reader, key, decodedtext)
	if err != nil {
		return "nil", err
	}
	return string(v15), err
}

func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	publicKey := privateKey.PublicKey
	fmt.Printf("privateKey:%s privateKey:%s\n", privateKey, publicKey)
}
