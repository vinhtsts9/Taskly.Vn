package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func GetHash(key string) string {
	hash := sha256.New()
	hash.Write([]byte(key))
	hashBytes := hash.Sum(nil)

	return hex.EncodeToString(hashBytes)
}

func GenerateSalt(length int) (string, error) {
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

func HassPassword(password string, salt string) string {
	//concatnate password and salt
	saltedPassword := password + salt
	hashPass := sha256.Sum256(([]byte(saltedPassword)))
	return hex.EncodeToString(hashPass[:])
}

func MatchingPassword(storeHash string, password string, salt string) bool {
	hashPassword := HassPassword(password, salt)
	return hashPassword == storeHash
}
