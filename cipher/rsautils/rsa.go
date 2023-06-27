package rsautils

import "crypto/rsa"

/**
 * @Author: lee
 * @Description:
 * @File: rsa
 * @Date: 2023-06-26 8:42 下午
 */

type RSACipher struct {
	pubKey []byte
	prvKey []byte
}

func NewRSACipher(pubKey, prvKey []byte) *RSACipher {
	ret := &RSACipher{}

	return ret
}

func (cipher *RSACipher) DecryptWithPubKey(plaintext []byte) ([]byte, error) {
	prvKey := rsa.DecryptPKCS1v15()
}
