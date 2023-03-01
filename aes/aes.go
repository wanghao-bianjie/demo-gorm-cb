package aes

import "strings"

func EncryptMock(s string) string {
	return "AES[" + s + "]"
}

func DecryptMock(s string) string {
	if strings.HasPrefix(s, "AES[") {
		s = strings.ReplaceAll(s, "AES[", "")
		s = s[:len(s)-1]
	}
	return s
}
