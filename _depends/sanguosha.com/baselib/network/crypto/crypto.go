package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

var (
	defalutAesKey = []byte{
		0x46, 0x72, 0x45, 0x6b, 0x55, 0x50, 0x37, 0x78,
		0x61, 0x4e, 0x3f, 0x26, 0x72, 0x65, 0x51, 0x3d,
		0x6a, 0x45, 0x66, 0x72, 0x61, 0x74, 0x68, 0x65,
		0x77, 0x35, 0x65, 0x47, 0x35, 0x51, 0x45, 0x63,
	}
	defalutAesIV = []byte{
		0x73, 0x65, 0x42, 0x37, 0x24, 0x46, 0x35, 0x53,
		0x23, 0x75, 0x66, 0x61, 0x6d, 0x55, 0x6d, 0x41,
	}
)

var block cipher.Block

func init() {
	var err error
	block, err = aes.NewCipher(defalutAesKey)
	if err != nil {
		panic(err)
	}
}

// Crypto ...
type Crypto interface {
	Encrypt(dst, src []byte)
	Decrypt(dst, src []byte) error
}

type aesCrypto struct {
	enc cipher.Stream
	dec cipher.Stream
}

// NewAesCrypto ...
func NewAesCrypto() Crypto {
	return &aesCrypto{
		enc: cipher.NewCTR(block, defalutAesIV),
		dec: cipher.NewCTR(block, defalutAesIV),
	}
}

// Encrypt ...
func (cpt *aesCrypto) Encrypt(dst, src []byte) {
	cpt.enc.XORKeyStream(dst, src)
}

// Decrypt ...
func (cpt *aesCrypto) Decrypt(dst, src []byte) error {
	cpt.dec.XORKeyStream(dst, src)
	return nil
}
