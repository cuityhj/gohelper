package pbe

import (
	"testing"

	ut "github.com/zdnscloud/cement/unittest"
)

func TestPbe(t *testing.T) {
	keyFactoryBase64 := "OGRhMGZmOTY1MTc5NWI3ZDYyMzFiNTVkMGFmY2RjYzU="
	workKey := "Linking@12345678Linking@12345678"
	password := "Linking@201907#$%^&*"
	decryptCtx, err := Encrypt(&EncryptContext{
		KeyFactoryBase64: keyFactoryBase64,
		WorkKey:          workKey,
		Password:         password,
		Iterator:         DefaultIterator,
	})

	ut.Assert(t, err == nil, "")
	decryptPassword, err := Decrypt(decryptCtx)
	ut.Assert(t, err == nil, "")
	ut.Assert(t, decryptPassword == password, decryptPassword)
}
