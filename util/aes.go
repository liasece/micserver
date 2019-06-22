package util

import (
	// "base/logger"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

const (
	KEY = "dfgjkue125js89ga"
	IV  = "ksgaiecv69dk36ka"
)

func AesEncrypt(origData []byte) ([]byte, error) {
	return AesEncryptByKey(origData, KEY, IV)
}

func AesEncryptByKey(origData []byte, aeskey string,
	aesiv string) ([]byte, error) {
	key := []byte(aeskey)
	iv := []byte(aesiv)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	origData = PKCS7Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(crypted []byte) ([]byte, error) {
	return AesDecryptByKey(crypted, KEY, IV)
}

func AesDecryptByKey(crypted []byte, aeskey string,
	aesiv string) ([]byte, error) {
	key := []byte(aeskey)
	iv := []byte(aesiv)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(crypted)%block.BlockSize() != 0 {
		return nil, errors.New("解密异常 密文长度不对")
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	if length == 0 {
		return []byte("")
	}
	unpadding := int(origData[length-1])
	if length < unpadding {
		return []byte("")
	}
	return origData[:(length - unpadding)]
}
