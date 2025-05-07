package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"
)

type Slug struct {
	Value string
}

func (s *Slug) GenerateSlug(length int, useAlphabetic, useNumeric bool) error {
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"

	characters := ""
	if useAlphabetic {
		characters += alphabet
	}
	if useNumeric {
		characters += digits
	}
	if characters == "" {
		return fmt.Errorf("at least one character type must be enabled")
	}

	result := make([]byte, length)
	max := big.NewInt(int64(len(characters)))

	for i := range length {
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return fmt.Errorf("failed to generate random number: %w", err)
		}
		result[i] = characters[num.Int64()]
	}

	randomPart := string(result)
	randomString := fmt.Sprintf("%s-%d", randomPart, time.Now().UnixMilli())

	s.Value = randomString

	return nil
}

func GenerateCSRFToken() string {
	b := make([]byte, 50)
	_, _ = rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
