package utilities

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"hash/fnv"
	"io"
)

// FNV32 hashes using fnv32 algorithm
func FNV32(text string) uint32 {
	algorithm := fnv.New32()
	return uint32Hasher(algorithm, text)
}

// FNV32a hashes using fnv32a algorithm
func FNV32a(text string) uint32 {
	algorithm := fnv.New32a()
	return uint32Hasher(algorithm, text)
}

// FNV64 hashes using fnv64 algorithm
func FNV64(text string) uint64 {
	algorithm := fnv.New64()
	return uint64Hasher(algorithm, text)
}

// FNV64a hashes using fnv64a algorithm
func FNV64a(text string) uint64 {
	algorithm := fnv.New64a()
	return uint64Hasher(algorithm, text)
}

// MD5 hashes using md5 algorithm
func MD5(text string) string {
	algorithm := md5.New()
	return stringHasher(algorithm, text)
}

// SHA1 hashes using sha1 algorithm
func SHA1(text string) string {
	algorithm := sha1.New()
	return stringHasher(algorithm, text)
}

// SHA256 hashes using sha256 algorithm
func SHA256(text string) string {
	algorithm := sha256.New()
	return stringHasher(algorithm, text)
}

// SHA512 hashes using sha512 algorithm
func SHA512(text string) string {
	algorithm := sha512.New()
	return stringHasher(algorithm, text)
}

func stringHasher(algorithm hash.Hash, text string) string {
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func uint32Hasher(algorithm hash.Hash32, text string) uint32 {
	_, err := algorithm.Write([]byte(text))
	if err != nil {
		return 0
	}
	return algorithm.Sum32()
}

func uint64Hasher(algorithm hash.Hash64, text string) uint64 {
	_, err := algorithm.Write([]byte(text))
	if err != nil {
		return 0
	}
	return algorithm.Sum64()
}

func EncryptAES(key string, message string) (string, error) {
	byteMsg := []byte(message)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("could not create new cipher: %v", err)
	}

	cipherText := make([]byte, aes.BlockSize+len(byteMsg))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("could not encrypt: %v", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], byteMsg)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptAES(key string, message string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", fmt.Errorf("could not base64 decode: %v", err)
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("could not create new cipher: %v", err)
	}

	if len(cipherText) < aes.BlockSize {
		return "", fmt.Errorf("invalid ciphertext block size")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
