package sm2

import (
	"bytes"
	"os"
	"testing"

	ut "github.com/cuityhj/cement/unittest"
)

func TestSM2(t *testing.T) {
	privKey, pubKey, err := GenerateSM2KeyPair()
	ut.Assert(t, err == nil, "generate sm2 key failed")

	plaintext := []byte("this is a test plaintext for verifying SM2 encryption and decryption")

	sign, err := privKey.Sign(plaintext)
	ut.Assert(t, err == nil, "sm2 sign failed")
	succeed, err := pubKey.Verify(plaintext, sign)
	ut.Assert(t, err == nil && succeed, "sm2 verify failed")
	ciphertext, err := pubKey.Encrypt(plaintext)
	ut.Assert(t, err == nil, "sm2 encrypt failed")
	decrypted, err := privKey.Decrypt(ciphertext)
	ut.Assert(t, err == nil, "sm2 decrypt failed")
	ut.Assert(t, bytes.Equal(decrypted, plaintext), "")

	privKeyFileName := "sm2_priv.pem"
	err = privKey.WriteToFile(privKeyFileName)
	ut.Assert(t, err == nil, "write sm2 private key failed")
	defer os.Remove(privKeyFileName)

	pubKeyFileName := "sm2_pub.pem"
	err = pubKey.WriteToFile(pubKeyFileName)
	ut.Assert(t, err == nil, "write sm2 public key failed")
	defer os.Remove(pubKeyFileName)

	privKey, err = LoadPrivateKeyFromFile(privKeyFileName)
	ut.Assert(t, err == nil, "load sm2 private key failed")

	pubKey, err = LoadPublicKeyFromFile(pubKeyFileName)
	ut.Assert(t, err == nil, "load sm2 public key failed")

	sign, err = privKey.Sign(plaintext)
	ut.Assert(t, err == nil, "sm2 sign failed")
	succeed, err = pubKey.Verify(plaintext, sign)
	ut.Assert(t, err == nil && succeed, "sm2 verify failed")
	ciphertext, err = pubKey.Encrypt(plaintext)
	ut.Assert(t, err == nil, "sm2 encrypt failed")
	decrypted, err = privKey.Decrypt(ciphertext)
	ut.Assert(t, err == nil, "sm2 decrypt failed")
	ut.Assert(t, bytes.Equal(decrypted, plaintext), "")
}
