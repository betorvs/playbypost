package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"math/rand"
)

func EncryptText(text, key string) (string, error) {
	// Convert key to bytes
	keyBytes := []byte(key)

	// new block cipher AES
	aesBlock, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// new mod CBC
	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return "", err
	}

	// new random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	// Text Encrypt
	textEncrypted := gcm.Seal(nil, nonce, []byte(text), nil)

	// Convert ciphertext to Base64
	ciphertextBase64 := base64.StdEncoding.EncodeToString(append(nonce, textEncrypted...))

	return ciphertextBase64, nil
}

func DecryptText(textEncrypted, key string) (string, error) {
	// fmt.Println("text Encrypted: ", textEncrypted)
	// Convert ciphertext from Base64 to bytes
	ciphertextBytes, err := base64.StdEncoding.DecodeString(textEncrypted)
	if err != nil {
		return "", err
	}

	// Convert key to bytes
	keyBytes := []byte(key)

	// new block cipher AES
	aesBlock, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// new mod to decrypt CBC
	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return "", err
	}

	// Split nonce from ciphertext
	nonce := ciphertextBytes[:gcm.NonceSize()]
	ciphertext := ciphertextBytes[gcm.NonceSize():]

	// text Decryption
	textDecrypted, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(textDecrypted), nil
}
