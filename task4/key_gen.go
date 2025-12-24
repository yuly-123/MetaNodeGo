package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

type KeyGenerator struct{}

// 对称加密，生成HMAC密钥
func (kg *KeyGenerator) GenerateHMACKey(bits int) (secretKey string, err error) {
	bytes := bits / 8
	key := make([]byte, bytes)
	_, err = rand.Read(key)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(key), err
}

// 非对称加密，生成RSA密钥对
func (kg *KeyGenerator) GenerateRSAKeyPair(bits int) (privateKey string, publicKey string, err error) {
	privKey, err := rsa.GenerateKey(rand.Reader, bits) // 生成私钥
	if err != nil {
		return "", "", err
	}

	// 编码私钥为PEM格式
	privKeyBytes := x509.MarshalPKCS1PrivateKey(privKey)
	privKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privKeyBytes,
	})

	// 编码公钥为PEM格式
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	return string(privKeyPEM), string(pubKeyPEM), nil
}

//func main() {
//	generator := &KeyGenerator{}
//
//	hmac256, err := generator.GenerateHMACKey(256)
//	fmt.Println("密钥（HMAC）:", hmac256, err)
//
//	privateKey, publicKey, err := generator.GenerateRSAKeyPair(2048)
//	fmt.Println("私钥 (RSA 2048):", privateKey, err)
//	fmt.Println("公钥 (RSA 2048):", publicKey, err)
//}
