package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
)

const secretKey = "x35k9f"

func GenerateUUID() string {
	return uuid.NewString()
}

func Decrypt(msg string) (string, error) {
	key := sha256.Sum256([]byte(secretKey))

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return "", err
	}

	nonce := key[len(key)-aesgcm.NonceSize():]

	encrypted, err := hex.DecodeString(msg)
	if err != nil {
		return "", err
	}

	decrypted, err := aesgcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

func Encrypt(src string) (string, error) {
	key := sha256.Sum256([]byte(secretKey))

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}

	nonce := key[len(key)-aesgcm.NonceSize():]

	dst := aesgcm.Seal(nil, nonce, []byte(src), nil)

	return hex.EncodeToString(dst), nil
}
