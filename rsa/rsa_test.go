package rsa

import (
	"testing"

	ut "github.com/cuityhj/cement/unittest"
)

func TestRsa(t *testing.T) {
	pubKeyPem := []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA6RUhL4k2IJiAAjni4gUb
H7T2OaRf3ZVGXYyaFJwwo7E6WbceXmJ54qaJ1z5Il2Fasopz0HgE26EIFMNW1aCn
ClIy6mQwc6RAdp1hPWOoswqR8p5nIbZoMKEaWd4H+U4l2br7sqfmPAJiz7MzBtzW
7hJ+yk2E0NnsBjbl35NzT1wyOk2i4XYQruVNmEVSZcSqfJKWoEgEdkrGSlZMdNz+
qqVJGVtDmR/SKLqoVLQmUEcNEiUUm1hiDY3rLqZlls14MsGoamOh9ms8jqkmXMdG
PFjdwEF6n3qsYcEMtiUPpwQXOjMjfBk3JN4uzFRO0h3v3lSnotK1sfcnWMx1o/RW
zQIDAQAB
-----END PUBLIC KEY-----`)
	privKeyPem := []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA6RUhL4k2IJiAAjni4gUbH7T2OaRf3ZVGXYyaFJwwo7E6Wbce
XmJ54qaJ1z5Il2Fasopz0HgE26EIFMNW1aCnClIy6mQwc6RAdp1hPWOoswqR8p5n
IbZoMKEaWd4H+U4l2br7sqfmPAJiz7MzBtzW7hJ+yk2E0NnsBjbl35NzT1wyOk2i
4XYQruVNmEVSZcSqfJKWoEgEdkrGSlZMdNz+qqVJGVtDmR/SKLqoVLQmUEcNEiUU
m1hiDY3rLqZlls14MsGoamOh9ms8jqkmXMdGPFjdwEF6n3qsYcEMtiUPpwQXOjMj
fBk3JN4uzFRO0h3v3lSnotK1sfcnWMx1o/RWzQIDAQABAoIBABK1f/rtVBtsskW7
eMM0rrcuUt9QmuNR4pqKuSlzWhIhnTHrWXQxCmCPhpvw42nHRjzzkRVINPeeJuRn
w5YVNaNF8tFOpnb51bTPsCj50WZBsHJdlaCmoBlbLR7OjCxOQgqGkEKIaJojhOPw
GeXwnw3wDR5W95S+EUx0ZqI6FzTk5eLIK0ZIZKGzjQUxb18irv7CI8IxCkqotBHr
56dQw27J1SM/xiaLnGpckeiRsd1CHOXJUqSKEIKQ+//TApYO3vcaju8Qv0LKq1Y0
Y6bpFV+xDAl3RX/cFWO5VcKI2nw7/s4eOA5SD90QzGnnovgegAslupW6VMZzr3Qd
wB1n27ECgYEA/da0NaEgr2us4t84yxEO1x3EBgnvj+TQ1XXMPPSScWDTnKKchUSZ
H3tKlZFOesflXlccSRHB3a2fIzOern8cik2nJB5Sitwh9xm4JM3ihXLdpSDS0F0m
MCWZ0MjPacbFOS0b3CmO3FCm+QvRWkouUTYOO5U96USbP06nRtbpD3kCgYEA6xEu
5bLPC6gwxAjo4PF7AhdWHoBb4xNpNpvxkQLcwKa7K3+t6REEAtqM8UyzAFVCszA1
OrGtb4o495V3ek6isZnLQ8MK94q0bKmJbjYqpVnltmVP0DxEafJLOOqOH0owdfO8
Sq2ovKN4X3KY3EJfGjWh5nJSLTUOVvSy1YriyPUCgYEA1ejlHHyYSrv2iYmLFrVd
SDKxSlV9KEmvIvOOFaAU+K6cJVdzh2rzjvAbPkehVx61T/cgwwLP2LvDa6rIgkxk
BLjDrVBQRuyTQuTNpVZLGiJeXhV3ElgtIk3NfYB8Katz8GbvH212Ent0+lLXLbtt
pMpk3Bk8fyNtoL/rf2sEJCECgYEA4iorNAeBG5ccLFDiNyM/lbh8PGaFggo4Dbmm
hz34xUbmCKkU24xqjpBWUQfZpbVismLso+c1ln6n5tYhGUU1ValCH5U2JQuIIpBP
0QE+sM64rG/3hcOmk0TyyPUr/sDztVjnzfYdxjmF5Feu6STWubHmboGJvUMx48oV
kk3Je00CgYEAsW+C7NLQhtrhSc6PY3DLLuJzDBFCHYDvqPOJjUuMh5OsPnlQvHYl
ezHsvPGH6LyH3Gv8zxaxpb+bNtDS07dZysvMxjahdPBYM+mf++/Kju14T/LF38GE
4xTCuQV92X8Mkb6CmWK2v9LZyGztsoL9CWD2h1EhXve87fzo2yjzlTk=
-----END RSA PRIVATE KEY-----
`)

	pubKey, err := ParsePublicKey(pubKeyPem)
	ut.Assert(t, err == nil, "parse public key failed")

	privKey, err := ParsePrivateKey(privKeyPem)
	ut.Assert(t, err == nil, "parse private key failed")

	password := "admin@123"
	encryptedPassword, err := EncryptWithPublicKey(password, pubKey)
	ut.Assert(t, err == nil, "encrypt password with public key failed")
	decryptedPassword, err := DecryptWithPrivateKey(encryptedPassword, privKey)
	ut.Assert(t, err == nil, "decrypt password with private key failed")
	ut.Assert(t, string(decryptedPassword) == password, string(decryptedPassword))
}
