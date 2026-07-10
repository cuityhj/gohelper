package sm2

import (
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

const (
	SM2_UNCOMPRESSED_C1_LENGTH = 65
	SM2_COMPRESSED_C1_LENGTH   = 33
	SM2_C3_LENGTH              = 32
)

func GenerateSM2KeyPair() (*sm2.PrivateKey, error) {
	return sm2.GenerateKey(rand.Reader)
}

func WritePrivateKeyToFile(privKey *sm2.PrivateKey, fileName string) error {
	privKeyPEM, err := x509.WritePrivateKeyToPem(privKey, nil)
	if err != nil {
		return fmt.Errorf("write private key to pem failed: %s", err.Error())
	}

	if err := os.WriteFile(fileName, privKeyPEM, 0600); err != nil {
		return fmt.Errorf("write private key pem to file failed: %s", err.Error())
	}

	return nil
}

func WritePublicKeyToFile(pubKey *sm2.PublicKey, fileName string) error {
	pubKeyPEM, err := x509.WritePublicKeyToPem(pubKey)
	if err != nil {
		return fmt.Errorf("write public key to pem failed: %s", err.Error())
	}

	if err := os.WriteFile(fileName, pubKeyPEM, 0644); err != nil {
		return fmt.Errorf("write public key pem to file failed: %s", err.Error())
	}

	return nil
}

func LoadPublicKeyFromFile(fileName string) (*sm2.PublicKey, error) {
	pubKeyPEM, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read sm2 public key file failed: %s", err.Error())
	}

	return ParsePublicKey(pubKeyPEM)
}

func ParsePublicKey(pubKeyPEM []byte) (*sm2.PublicKey, error) {
	pubKey, err := x509.ReadPublicKeyFromPem(pubKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("gmsm x509 parse public key failed: %s", err.Error())
	}

	return pubKey, nil
}

func EncryptWithPublicPEM(plaintext string, pubKeyPEM []byte) (string, error) {
	pubKey, err := ParsePublicKey(pubKeyPEM)
	if err != nil {
		return "", err
	}

	return Encrypt(pubKey, []byte(plaintext))
}

func Encrypt(pubKey *sm2.PublicKey, plaintext []byte) (string, error) {
	rawCipher, err := sm2.Encrypt(pubKey, plaintext, rand.Reader, sm2.C1C3C2)
	if err != nil {
		return "", fmt.Errorf("sm2 encrypt plaintext failed: %s", err.Error())
	}

	compressedCipher, err := compressSM2Ciphertext(rawCipher)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(compressedCipher), nil
}

func compressSM2Ciphertext(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < SM2_UNCOMPRESSED_C1_LENGTH+SM2_C3_LENGTH {
		return nil, fmt.Errorf("invalid uncompressed sm2 ciphertext length %d", len(ciphertext))
	}

	curve := sm2.P256Sm2()
	x, y := elliptic.Unmarshal(curve, ciphertext[:SM2_UNCOMPRESSED_C1_LENGTH])
	if x == nil {
		return nil, fmt.Errorf("invalid uncompressed C1 point in ciphertext")
	}

	c1Compressed := elliptic.MarshalCompressed(curve, x, y)
	result := make([]byte, len(c1Compressed)+len(ciphertext[SM2_UNCOMPRESSED_C1_LENGTH:]))
	copy(result[:len(c1Compressed)], c1Compressed)
	copy(result[len(c1Compressed):], ciphertext[SM2_UNCOMPRESSED_C1_LENGTH:])
	return result, nil
}

func LoadPrivateKeyFromFile(fileName string) (*sm2.PrivateKey, error) {
	privKeyPEM, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read sm2 private key file failed: %s", err.Error())
	}

	return ParsePrivateKey(privKeyPEM)
}

func ParsePrivateKey(privKeyPEM []byte) (*sm2.PrivateKey, error) {
	privKey, err := x509.ReadPrivateKeyFromPem(privKeyPEM, nil)
	if err != nil {
		return nil, fmt.Errorf("gmsm x509 parse private key failed: %s", err.Error())
	}

	return privKey, nil
}

func DecryptWithPrivatePEM(encryptedData string, privKeyPEM []byte) ([]byte, error) {
	privKey, err := ParsePrivateKey(privKeyPEM)
	if err != nil {
		return nil, err
	}

	return Decrypt(privKey, encryptedData)
}

func Decrypt(privKey *sm2.PrivateKey, encryptedData string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("base64 decode encrypted data failed: %s", err.Error())
	}

	rawCipher, err := decompressSM2Ciphertext(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("decompress sm2 ciphertext failed: %s", err.Error())
	}

	plaintext, err := sm2.Decrypt(privKey, rawCipher, sm2.C1C3C2)
	if err != nil {
		return nil, fmt.Errorf("sm2 decrypt ciphertext failed: %s", err.Error())
	}

	return plaintext, nil
}

func decompressSM2Ciphertext(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < SM2_COMPRESSED_C1_LENGTH+SM2_C3_LENGTH {
		return nil, fmt.Errorf("invalid compressed sm2 ciphertext length %d", len(ciphertext))
	}

	curve := sm2.P256Sm2()
	x, y := elliptic.UnmarshalCompressed(curve, ciphertext[:SM2_COMPRESSED_C1_LENGTH])
	if x == nil {
		return nil, fmt.Errorf("invalid compressed C1 point in ciphertext")
	}

	c1Raw := elliptic.Marshal(curve, x, y)
	result := make([]byte, len(c1Raw)+len(ciphertext[SM2_COMPRESSED_C1_LENGTH:]))
	copy(result[:len(c1Raw)], c1Raw)
	copy(result[len(c1Raw):], ciphertext[SM2_COMPRESSED_C1_LENGTH:])
	return result, nil
}
