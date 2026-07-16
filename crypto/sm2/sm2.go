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

type PrivateKey struct {
	*sm2.PrivateKey
}

type PublicKey struct {
	*sm2.PublicKey
}

func GenerateSM2KeyPair() (*PrivateKey, *PublicKey, error) {
	if privKey, err := sm2.GenerateKey(rand.Reader); err != nil {
		return nil, nil, err
	} else {
		return &PrivateKey{PrivateKey: privKey}, &PublicKey{PublicKey: &privKey.PublicKey}, nil
	}
}

func (priv *PrivateKey) Sign(msg []byte) (string, error) {
	if len(msg) == 0 {
		return "", fmt.Errorf("message cannot be empty")
	}

	sign, err := priv.PrivateKey.Sign(rand.Reader, msg, nil)
	if err != nil {
		return "", fmt.Errorf("sm2 private key sign failed: %w", err)
	}

	return base64.StdEncoding.EncodeToString(sign), nil
}

func (pub *PublicKey) Verify(msg []byte, encryptedSign string) (bool, error) {
	if len(msg) == 0 || len(encryptedSign) == 0 {
		return false, fmt.Errorf("sign or message cannot be empty")
	}

	sign, err := base64.StdEncoding.DecodeString(encryptedSign)
	if err != nil {
		return false, fmt.Errorf("base64 decode sign failed: %w", err)
	}

	return pub.PublicKey.Verify(msg, sign), nil
}

func (pub *PublicKey) WriteToFile(fileName string) error {
	pubKeyPEM, err := pub.ToPEM()
	if err != nil {
		return fmt.Errorf("write public key to pem failed: %w", err)
	}

	if err := os.WriteFile(fileName, pubKeyPEM, 0644); err != nil {
		return fmt.Errorf("write public key pem to file failed: %w", err)
	}

	return nil
}

func (pub *PublicKey) ToPEM() ([]byte, error) {
	return x509.WritePublicKeyToPem(pub.PublicKey)
}

func LoadPublicKeyFromFile(fileName string) (*PublicKey, error) {
	pubKeyPEM, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read sm2 public key file failed: %w", err)
	}

	return ParsePublicKey(pubKeyPEM)
}

func ParsePublicKey(pubKeyPEM []byte) (*PublicKey, error) {
	if len(pubKeyPEM) == 0 {
		return nil, fmt.Errorf("public pem cannot be empty")
	}

	pubKey, err := x509.ReadPublicKeyFromPem(pubKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("gmsm x509 parse public key failed: %w", err)
	}

	return &PublicKey{pubKey}, nil
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

	rawCipher, err := sm2.Encrypt(pub.PublicKey, plaintext, rand.Reader, sm2.C1C3C2)
	if err != nil {
		return "", fmt.Errorf("sm2 encrypt plaintext failed: %w", err)
	}

	compressedCipher, err := compressSM2Ciphertext(rawCipher)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(compressedCipher), nil
}

func EncryptAsn1WithPublicPEM(plaintext string, pubKeyPEM []byte) (string, error) {
	if len(plaintext) == 0 || len(pubKeyPEM) == 0 {
		return "", fmt.Errorf("plaintext or public pem cannot be empty")
	}

	pubKey, err := ParsePublicKey(pubKeyPEM)
	if err != nil {
		return "", err
	}

	return pubKey.EncryptAsn1([]byte(plaintext))
}

func (pub *PublicKey) EncryptAsn1(plaintext []byte) (string, error) {
	ciphertext, err := pub.PublicKey.EncryptAsn1(plaintext, rand.Reader)
	if err != nil {
		return "", fmt.Errorf("sm2 encrypt asn.1 plaintext failed: %w", err)
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
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
	return &PublicKey{&priv.PrivateKey.PublicKey}
}

func (priv *PrivateKey) WriteToFile(fileName string) error {
	privKeyPEM, err := priv.ToPEM()
	if err != nil {
		return fmt.Errorf("write private key to pem failed: %w", err)
	}

	if err := os.WriteFile(fileName, privKeyPEM, 0600); err != nil {
		return fmt.Errorf("write private key pem to file failed: %w", err)
	}

	return nil
}

func (priv *PrivateKey) ToPEM() ([]byte, error) {
	return x509.WritePrivateKeyToPem(priv.PrivateKey, nil)
}

func (priv *PrivateKey) WriteEncryptedToFile(fileName string, password []byte) error {
	privKeyPEM, err := priv.ToEncryptedPEM(password)
	if err != nil {
		return fmt.Errorf("write private key to pem failed: %w", err)
	}

	if err := os.WriteFile(fileName, privKeyPEM, 0600); err != nil {
		return fmt.Errorf("write private key pem to file failed: %w", err)
	}

	return nil
}

func (priv *PrivateKey) ToEncryptedPEM(password []byte) ([]byte, error) {
	if len(password) == 0 {
		return nil, fmt.Errorf("password cannot be empty")
	}

	return x509.WritePrivateKeyToPem(priv.PrivateKey, password)
}

func LoadPrivateKeyFromFile(fileName string) (*PrivateKey, error) {
	privKeyPEM, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read sm2 private key file failed: %w", err)
	}

	return ParsePrivateKey(privKeyPEM)
}

func ParsePrivateKey(privKeyPEM []byte) (*PrivateKey, error) {
	if len(privKeyPEM) == 0 {
		return nil, fmt.Errorf("private pem cannot be empty")
	}

	privKey, err := x509.ReadPrivateKeyFromPem(privKeyPEM, nil)
	if err != nil {
		return nil, fmt.Errorf("gmsm x509 parse private key failed: %w", err)
	}

	return &PrivateKey{privKey}, nil
}

func LoadEncryptedPrivateKeyFromFile(fileName string, password []byte) (*PrivateKey, error) {
	if len(password) == 0 {
		return nil, fmt.Errorf("password cannot be empty")
	}

	privKeyPEM, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read sm2 private key file failed: %w", err)
	}

	return ParseEncryptedPrivateKey(privKeyPEM, password)
}

func ParseEncryptedPrivateKey(privKeyPEM, password []byte) (*PrivateKey, error) {
	if len(privKeyPEM) == 0 || len(password) == 0 {
		return nil, fmt.Errorf("private pem or password cannot be empty")
	}

	privKey, err := x509.ReadPrivateKeyFromPem(privKeyPEM, password)
	if err != nil {
		return nil, fmt.Errorf("gmsm x509 parse private key failed: %w", err)
	}

	return &PrivateKey{privKey}, nil
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
		return nil, fmt.Errorf("base64 decode encrypted data failed: %w", err)
	}

	rawCipher := ciphertext
	if len(ciphertext) > 0 && (ciphertext[0] == 0x02 || ciphertext[0] == 0x03) {
		rawCipher, err = decompressSM2Ciphertext(ciphertext)
		if err != nil {
			return nil, fmt.Errorf("decompress sm2 ciphertext failed: %w", err)
		}
	}

	plaintext, err := sm2.Decrypt(priv.PrivateKey, rawCipher, sm2.C1C3C2)
	if err != nil {
		return nil, fmt.Errorf("sm2 decrypt ciphertext failed: %w", err)
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

func DecryptAsn1WithPrivatePEM(encryptedData string, privKeyPEM []byte) ([]byte, error) {
	if len(encryptedData) == 0 || len(privKeyPEM) == 0 {
		return nil, fmt.Errorf("ciphertext or private pem cannot be empty")
	}

	privKey, err := ParsePrivateKey(privKeyPEM)
	if err != nil {
		return nil, err
	}

	return privKey.DecryptAsn1(encryptedData)
}

func DecryptAsn1WithEncryptedPrivatePEM(encryptedData string, privKeyPEM, password []byte) ([]byte, error) {
	if len(encryptedData) == 0 || len(privKeyPEM) == 0 || len(password) == 0 {
		return nil, fmt.Errorf("ciphertext or private pem or password cannot be empty")
	}

	privKey, err := ParseEncryptedPrivateKey(privKeyPEM, password)
	if err != nil {
		return nil, err
	}

	return privKey.DecryptAsn1(encryptedData)
}

func (priv *PrivateKey) DecryptAsn1(encryptedData string) ([]byte, error) {
	if len(encryptedData) == 0 {
		return nil, fmt.Errorf("ciphertext cannot be empty")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("base64 decode encrypted data failed: %w", err)
	}

	plaintext, err := priv.PrivateKey.DecryptAsn1(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("sm2 decrypt asn.1 ciphertext failed: %w", err)
	}

	return plaintext, nil
}
