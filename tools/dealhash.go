package tools

import (
	"crypto/sha256"
	"fmt"
)

func CodeHash(code string) string{
	h := sha256.New()
	h.Write([]byte(code))
	pass_hash := fmt.Sprintf("%x", h.Sum(nil))
	return pass_hash
}