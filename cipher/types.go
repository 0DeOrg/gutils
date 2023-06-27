package cipher

/**
 * @Author: lee
 * @Description:
 * @File: types
 * @Date: 2023-06-26 8:43 下午
 */

type ICipher interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}
