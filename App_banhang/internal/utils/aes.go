package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
)

func EncryptAES(plainText []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return "", err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, plainText, nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func DecryptAES(cipherBase64 string, key []byte) ([]byte, error) {
	cipherText, err := base64.URLEncoding.DecodeString(cipherBase64)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := cipherText[:nonceSize], cipherText[nonceSize:]
	return aesGCM.Open(nil, nonce, ciphertext, nil)
}
