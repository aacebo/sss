package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

func AESEncrypt(value []byte, secret string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(secret))

	if err != nil {
		return nil, err
	}

	data := []byte(value)
	text := make([]byte, aes.BlockSize+len(data))
	iv := text[:aes.BlockSize]

	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(text[aes.BlockSize:], data)

	return []byte(base64.URLEncoding.EncodeToString(text)), nil
}

func AESDecrypt(value []byte, secret string) ([]byte, error) {
	text, err := base64.URLEncoding.DecodeString(string(value))

	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(secret))

	if err != nil {
		return nil, err
	}

	if len(text) < aes.BlockSize {
		return nil, errors.New("invalid cypher")
	}

	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(text, text)
	return text, nil
}
