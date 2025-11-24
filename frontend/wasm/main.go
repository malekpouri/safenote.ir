package main // Force rebuild

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"syscall/js"
)

func deriveKey(shortKey, password string) []byte {
	hash := sha256.Sum256([]byte(shortKey + password))
	return hash[:]
}

func encrypt(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return map[string]interface{}{
			"error": "Invalid arguments. Expected (text, shortKey, [password])",
		}
	}

	text := args[0].String()
	shortKey := args[1].String()
	password := ""
	if len(args) > 2 {
		password = args[2].String()
	}

	key := deriveKey(shortKey, password)

	block, err := aes.NewCipher(key)
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(text), nil)
	return map[string]interface{}{
		"encrypted": base64.StdEncoding.EncodeToString(ciphertext),
	}
}

func decrypt(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return map[string]interface{}{
			"error": "Invalid arguments. Expected (encryptedBase64, shortKey, [password])",
		}
	}

	encryptedBase64 := args[0].String()
	shortKey := args[1].String()
	password := ""
	if len(args) > 2 {
		password = args[2].String()
	}

	key := deriveKey(shortKey, password)

	data, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return map[string]interface{}{
			"error": "Invalid base64 data",
		}
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return map[string]interface{}{
			"error": "Ciphertext too short",
		}
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	return map[string]interface{}{
		"decrypted": string(plaintext),
	}
}

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("encryptNote", js.FuncOf(encrypt))
	js.Global().Set("decryptNote", js.FuncOf(decrypt))
	<-c
}
