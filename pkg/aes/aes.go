package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

func normalizeKey(key string) []byte {
	k := []byte(key)
	if len(k) >= aes.BlockSize {
		return k[:aes.BlockSize]
	}

	padded := make([]byte, aes.BlockSize)
	copy(padded, k)

	return padded
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(data, padText...)
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("invalid padding")
	}

	padding := int(data[length-1])
	if padding == 0 || padding > length || padding > aes.BlockSize {
		return nil, errors.New("invalid padding")
	}

	for _, b := range data[length-padding:] {
		if int(b) != padding {
			return nil, errors.New("invalid padding")
		}
	}

	return data[:length-padding], nil
}

func Encrypt(plainText, key string) (string, error) {
	keyBytes := normalizeKey(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	plainBytes := pkcs7Pad([]byte(plainText), block.BlockSize())
	encrypted := make([]byte, len(plainBytes))
	mode := cipher.NewCBCEncrypter(block, keyBytes)
	mode.CryptBlocks(encrypted, plainBytes)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func Decrypt(cipherText, key string) (string, error) {
	keyBytes := normalizeKey(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	raw, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	if len(raw)%block.BlockSize() != 0 {
		return "", errors.New("invalid ciphertext length")
	}

	decrypted := make([]byte, len(raw))
	mode := cipher.NewCBCDecrypter(block, keyBytes)
	mode.CryptBlocks(decrypted, raw)

	plainBytes, err := pkcs7Unpad(decrypted)
	if err != nil {
		return "", err
	}

	return string(plainBytes), nil
}
