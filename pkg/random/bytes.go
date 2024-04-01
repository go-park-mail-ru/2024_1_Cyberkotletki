package random

import "crypto/rand"

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func bytesWithCharset(length int, charset string) ([]byte, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	for i := range bytes {
		bytes[i] = charset[int(bytes[i])%len(charset)]
	}
	return bytes, nil
}

// Bytes генерирует случайную строку длиной length из символов английского алфавита и цифр
func Bytes(length int) ([]byte, error) {
	return bytesWithCharset(length, charset)
}
