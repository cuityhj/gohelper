package sm2

import (
	"bytes"
	"os"
	"testing"

	ut "github.com/cuityhj/cement/unittest"
)

func TestSM2(t *testing.T) {
	privKey, err := GenerateSM2KeyPair()
	ut.Assert(t, err == nil, "generate sm2 key failed")

	pubKey := &privKey.PublicKey

	plaintext := []byte("this is a test plaintext for verifying SM2 encryption and decryption")
	ciphertext, err := Encrypt(pubKey, plaintext)
	ut.Assert(t, err == nil, "sm2 encrypt failed")
	decrypted, err := Decrypt(privKey, ciphertext)
	ut.Assert(t, err == nil, "sm2 decrypt failed")
	ut.Assert(t, bytes.Equal(decrypted, plaintext), "")

	privKeyFileName := "sm2_priv.pem"
	err = WritePrivateKeyToFile(privKey, privKeyFileName)
	ut.Assert(t, err == nil, "write sm2 private key failed")
	defer os.Remove(privKeyFileName)

	pubKeyFileName := "sm2_pub.pem"
	err = WritePublicKeyToFile(pubKey, pubKeyFileName)
	ut.Assert(t, err == nil, "write sm2 public key failed")
	defer os.Remove(pubKeyFileName)

	privKey, err = LoadPrivateKeyFromFile(privKeyFileName)
	ut.Assert(t, err == nil, "load sm2 private key failed")

	pubKey, err = LoadPublicKeyFromFile(pubKeyFileName)
	ut.Assert(t, err == nil, "load sm2 public key failed")

	ciphertext, err = Encrypt(pubKey, plaintext)
	ut.Assert(t, err == nil, "sm2 encrypt failed")
	decrypted, err = Decrypt(privKey, ciphertext)
	ut.Assert(t, err == nil, "sm2 decrypt failed")
	ut.Assert(t, bytes.Equal(decrypted, plaintext), "")
}
