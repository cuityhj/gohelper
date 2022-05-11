package pbe

import (
	"crypto/sha512"
	"encoding/base64"

	"golang.org/x/crypto/pbkdf2"
)

var (
	rootVector = []byte("lg!vYve07'E=,u).AGZHU;u6-Z;j:[$|4L4Jbni6ygp8g=ed]Jt&rHjr0Fg2,oQX5{G]PMz.t+7UQMP!>x1cC,iMlC]ZI=]6Y:X6&ZyhyRMJl'~$:I:Dg~U+KMN,-7MmrQF*5%njJ:,8ebyP;;zX9Hq(n97cx&T:)B;&%aROq)b4l|Ah<S(S4DFKvHPJVupLDIXmM(b>C8^fOD|<xpWC&FSKdI}Wa@VWTW##RNJ@CFc-uL?+C/>-KI:xl$:.t,3xm1(V7_<F#eil;H;,Qdxms!&t;5:>0L'*f#uSb52zWhD$D9`9VMO~,>~&eO8XExc,eL{x!8b,elq6Y@?[kq,NC{yQTmjat<`z<&)u;{x#q@LF.;=YPpQe9Z*6;&N*$ua}eMeIWS/;=6nM~X>jct9>E]7u@~XTJcXs9v%:%u{z@N7vmCKyXI9xu0RU^Z'kT3'G+hKBbakjWu~kKh,ji:tORF^Hh1fbto;VS~7NM$I[M,,Eh^([Pz(#ULUdrc-,94Ur")
)

type DecryptContext struct {
	KeyFactoryBase64 string
	EncryptWorkKey   string
	EncryptPassword  string
	Iterator         int
}

func Decrypt(ctx *DecryptContext) (string, error) {
	keyFactory, err := base64.StdEncoding.DecodeString(ctx.KeyFactoryBase64)
	if err != nil {
		return "", err
	}

	rootKey := string(pbkdf2.Key(rootVector, keyFactory, ctx.Iterator, sha512.BlockSize, sha512.New))
	decryptWorkKey, err := CBCDecrypt(rootKey[:32], ctx.EncryptWorkKey)
	if err != nil {
		return "", err
	}

	return CBCDecrypt(decryptWorkKey, ctx.EncryptPassword)
}

type EncryptContext struct {
	KeyFactoryBase64 string
	WorkKey          string
	Password         string
	Iterator         int
}

func Encrypt(ctx *EncryptContext) (*DecryptContext, error) {
	keyFactory, err := base64.StdEncoding.DecodeString(ctx.KeyFactoryBase64)
	if err != nil {
		return nil, err
	}

	rootKey := string(pbkdf2.Key(rootVector, keyFactory, ctx.Iterator, sha512.BlockSize, sha512.New))
	encryptWorkKey, err := CBCEncrypt(rootKey[:32], ctx.WorkKey)
	if err != nil {
		return nil, err
	}

	encryptPassword, err := CBCEncrypt(ctx.WorkKey, ctx.Password)
	if err != nil {
		return nil, err
	}

	return &DecryptContext{
		KeyFactoryBase64: ctx.KeyFactoryBase64,
		EncryptWorkKey:   encryptWorkKey,
		EncryptPassword:  encryptPassword,
		Iterator:         ctx.Iterator,
	}, nil
}
