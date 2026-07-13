package pbe

import (
	"testing"

	ut "github.com/cuityhj/cement/unittest"
)

func TestPbe(t *testing.T) {
	keyFactoryBase64 := "OGRhMGZmOTY1MTc5NWI3ZDYyMzFiNTVkMGFmY2RjYzU="
	workKey := "Linking@12345678Linking@12345678"
	plaintext := "Linking@202004#$%^&*"
	decryptCtx, err := Encrypt(&EncryptContext{
		KeyFactoryBase64: keyFactoryBase64,
		WorkKey:          workKey,
		Plaintext:        plaintext,
		Iterator:         DEFAULT_ITERATOR,
	})

	ut.Assert(t, err == nil, "")
	decrypted, err := Decrypt(decryptCtx)
	ut.Assert(t, err == nil, "")
	ut.Assert(t, string(decrypted) == plaintext, string(decrypted))
}
