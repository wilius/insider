package types

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

var aesBlock cipher.Block

func init() {
	var err error
	aesBlock, err = aes.NewCipher([]byte("zndnMaT0tCi9Ia73BUSwP93M8BWOzMHZ"))
	if err != nil {
		panic(fmt.Sprintf("Failed to create AES cipher: %v", err))
	}
}

func Encrypt(text *string) (*string, error) {
	plaintext := []byte(*text)
	blockSize := aesBlock.BlockSize()
	paddedPlaintext := pad(plaintext, blockSize)

	ciphertext := make([]byte, len(paddedPlaintext))
	for i := 0; i < len(paddedPlaintext); i += blockSize {
		aesBlock.Encrypt(ciphertext[i:i+blockSize], paddedPlaintext[i:i+blockSize])
	}

	toString := hex.EncodeToString(ciphertext)
	return &toString, nil
}

func Decrypt(text *string) (*string, error) {
	ciphertext, err := hex.DecodeString(*text)
	if err != nil {
		return nil, err
	}

	blockSize := aesBlock.BlockSize()
	if len(ciphertext)%blockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	plaintext := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += blockSize {
		aesBlock.Decrypt(plaintext[i:i+blockSize], ciphertext[i:i+blockSize])
	}

	plaintext = unpad(plaintext)
	s := string(plaintext)
	return &s, nil
}

func pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func unpad(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
