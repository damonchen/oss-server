package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
)

var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

//GetRandomString ...
func GetRandomString(l int) string {
	b := make([]rune, l)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetFilename(v string) string {
	s := sha1.New()
	s.Write([]byte(v))
	sum := s.Sum(nil)
	return hex.EncodeToString(sum)
}
