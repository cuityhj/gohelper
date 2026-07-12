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
	// SM2 未压缩公钥 C1 长度：1字节标识 + 32字节X + 32字节Y = 65字节
	SM2_UNCOMPRESSED_C1_LENGTH = 65
	// SM2 压缩公钥 C1 长度：1字节标识 + 32字节X = 33字节
	SM2_COMPRESSED_C1_LENGTH = 33
	// SM2 C3 摘要长度：SM3 哈希输出 32 字节
	SM2_C3_LENGTH = 32
)

var sm2Curve = sm2.P256Sm2()

type (
	PrivateKey sm2.PrivateKey
	PublicKey  sm2.PublicKey
)

func GenerateSM2KeyPair() (*PrivateKey, *PublicKey, error) {
	if privKey, err := sm2.GenerateKey(rand.Reader); err != nil {
		return nil, nil, err
	} else {
		return (*PrivateKey)(privKey), (*PublicKey)(&privKey.PublicKey), nil
	}
}

func (pub *PublicKey) ToHex() string {
	return fmt.Sprintf("04%064x%064x", pub.X, pub.Y)
}

func (pub *PublicKey) WriteToFile(fileName string) error {
	pubKeyPEM, err := pub.ToPEM()
	if err != nil {
		return fmt.Errorf("write public key to pem failed: %s", err.Error())
	}

	if err := os.WriteFile(fileName, pubKeyPEM, 0644); err != nil {
		return fmt.Errorf("write public key pem to file failed: %s", err.Error())
	}

	return nil
}

func (pub *PublicKey) ToPEM() ([]byte, error) {
	return x509.WritePublicKeyToPem((*sm2.PublicKey)(pub))
}

func LoadPublicKeyFromFile(fileName string) (*PublicKey, error) {
	pubKeyPEM, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read sm2 public key file failed: %s", err.Error())
	}

	return ParsePublicKey(pubKeyPEM)
}

func ParsePublicKey(pubKeyPEM []byte) (*PublicKey, error) {
	if len(pubKeyPEM) == 0 {
		return nil, fmt.Errorf("public pem cannot be empty")
	}

	pubKey, err := x509.ReadPublicKeyFromPem(pubKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("gmsm x509 parse public key failed: %s", err.Error())
	}

	return (*PublicKey)(pubKey), nil
}

func EncryptWithPublicPEM(plaintext string, pubKeyPEM []byte) (string, error) {
	if len(plaintext) == 0 || len(pubKeyPEM) == 0 {
		return "", fmt.Errorf("plaintext or public pem cannot be empty")
	}

	pubKey, err := ParsePublicKey(pubKeyPEM)
	if err != nil {
		return "", err
	}

	return pubKey.Encrypt([]byte(plaintext))
}

func (pub *PublicKey) Encrypt(plaintext []byte) (string, error) {
	if len(plaintext) == 0 {
		return "", fmt.Errorf("plaintext cannot be empty")
	}

	rawCipher, err := sm2.Encrypt((*sm2.PublicKey)(pub), plaintext, rand.Reader, sm2.C1C3C2)
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

	x, y := elliptic.Unmarshal(sm2Curve, ciphertext[:SM2_UNCOMPRESSED_C1_LENGTH])
	if x == nil {
		return nil, fmt.Errorf("invalid uncompressed C1 point in ciphertext")
	}

	c1Compressed := elliptic.MarshalCompressed(sm2Curve, x, y)
	result := make([]byte, len(c1Compressed)+len(ciphertext[SM2_UNCOMPRESSED_C1_LENGTH:]))
	copy(result[:len(c1Compressed)], c1Compressed)
	copy(result[len(c1Compressed):], ciphertext[SM2_UNCOMPRESSED_C1_LENGTH:])
	return result, nil
}

func (priv *PrivateKey) Public() *PublicKey {
	return (*PublicKey)(&priv.PublicKey)
}

func (priv *PrivateKey) WriteToFile(fileName string) error {
	privKeyPEM, err := priv.ToPEM()
	if err != nil {
		return fmt.Errorf("write private key to pem failed: %s", err.Error())
	}

	if err := os.WriteFile(fileName, privKeyPEM, 0600); err != nil {
		return fmt.Errorf("write private key pem to file failed: %s", err.Error())
	}

	return nil
}

func (priv *PrivateKey) ToPEM() ([]byte, error) {
	return x509.WritePrivateKeyToPem((*sm2.PrivateKey)(priv), nil)
}

func (priv *PrivateKey) WriteEncryptedToFile(fileName string, password []byte) error {
	privKeyPEM, err := priv.ToEncryptedPEM(password)
	if err != nil {
		return fmt.Errorf("write private key to pem failed: %s", err.Error())
	}

	if err := os.WriteFile(fileName, privKeyPEM, 0600); err != nil {
		return fmt.Errorf("write private key pem to file failed: %s", err.Error())
	}

	return nil
}

func (priv *PrivateKey) ToEncryptedPEM(password []byte) ([]byte, error) {
	if len(password) == 0 {
		return nil, fmt.Errorf("password cannot be empty")
	}

	return x509.WritePrivateKeyToPem((*sm2.PrivateKey)(priv), password)
}

func LoadPrivateKeyFromFile(fileName string) (*PrivateKey, error) {
	privKeyPEM, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read sm2 private key file failed: %s", err.Error())
	}

	return ParsePrivateKey(privKeyPEM)
}

func ParsePrivateKey(privKeyPEM []byte) (*PrivateKey, error) {
	if len(privKeyPEM) == 0 {
		return nil, fmt.Errorf("private pem cannot be empty")
	}

	privKey, err := x509.ReadPrivateKeyFromPem(privKeyPEM, nil)
	if err != nil {
		return nil, fmt.Errorf("gmsm x509 parse private key failed: %s", err.Error())
	}

	return (*PrivateKey)(privKey), nil
}

func LoadEncryptedPrivateKeyFromFile(fileName string, password []byte) (*PrivateKey, error) {
	if len(password) == 0 {
		return nil, fmt.Errorf("password cannot be empty")
	}

	privKeyPEM, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read sm2 private key file failed: %s", err.Error())
	}

	return ParseEncryptedPrivateKey(privKeyPEM, password)
}

func ParseEncryptedPrivateKey(privKeyPEM, password []byte) (*PrivateKey, error) {
	if len(privKeyPEM) == 0 || len(password) == 0 {
		return nil, fmt.Errorf("private pem or password cannot be empty")
	}

	privKey, err := x509.ReadPrivateKeyFromPem(privKeyPEM, password)
	if err != nil {
		return nil, fmt.Errorf("gmsm x509 parse private key failed: %s", err.Error())
	}

	return (*PrivateKey)(privKey), nil
}

func DecryptWithPrivatePEM(encryptedData string, privKeyPEM []byte) ([]byte, error) {
	if len(encryptedData) == 0 || len(privKeyPEM) == 0 {
		return nil, fmt.Errorf("ciphertext or private pem cannot be empty")
	}

	privKey, err := ParsePrivateKey(privKeyPEM)
	if err != nil {
		return nil, err
	}

	return privKey.Decrypt(encryptedData)
}

func DecryptWithEncryptedPrivatePEM(encryptedData string, privKeyPEM, password []byte) ([]byte, error) {
	if len(encryptedData) == 0 || len(privKeyPEM) == 0 || len(password) == 0 {
		return nil, fmt.Errorf("ciphertext or private pem or password cannot be empty")
	}

	privKey, err := ParseEncryptedPrivateKey(privKeyPEM, password)
	if err != nil {
		return nil, err
	}

	return privKey.Decrypt(encryptedData)
}

func (priv *PrivateKey) Decrypt(encryptedData string) ([]byte, error) {
	if len(encryptedData) == 0 {
		return nil, fmt.Errorf("ciphertext cannot be empty")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("base64 decode encrypted data failed: %s", err.Error())
	}

	rawCipher := ciphertext
	if len(ciphertext) > 0 && (ciphertext[0] == 0x02 || ciphertext[0] == 0x03) {
		rawCipher, err = decompressSM2Ciphertext(ciphertext)
		if err != nil {
			return nil, fmt.Errorf("decompress sm2 ciphertext failed: %s", err.Error())
		}
	}

	plaintext, err := sm2.Decrypt((*sm2.PrivateKey)(priv), rawCipher, sm2.C1C3C2)
	if err != nil {
		return nil, fmt.Errorf("sm2 decrypt ciphertext failed: %s", err.Error())
	}

	return plaintext, nil
}

func decompressSM2Ciphertext(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < SM2_COMPRESSED_C1_LENGTH+SM2_C3_LENGTH {
		return nil, fmt.Errorf("invalid compressed sm2 ciphertext length %d", len(ciphertext))
	}

	x, y := elliptic.UnmarshalCompressed(sm2Curve, ciphertext[:SM2_COMPRESSED_C1_LENGTH])
	if x == nil {
		return nil, fmt.Errorf("invalid compressed C1 point in ciphertext")
	}

	c1Raw := elliptic.Marshal(sm2Curve, x, y)
	result := make([]byte, len(c1Raw)+len(ciphertext[SM2_COMPRESSED_C1_LENGTH:]))
	copy(result[:len(c1Raw)], c1Raw)
	copy(result[len(c1Raw):], ciphertext[SM2_COMPRESSED_C1_LENGTH:])
	return result, nil
}
