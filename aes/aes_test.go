package aes

import "testing"

func TestAESMock(t *testing.T) {
	a := EncryptMock("qwe")
	t.Log(a)
	t.Log(DecryptMock(a))
}
