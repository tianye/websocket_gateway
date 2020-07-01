package encry

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"math/rand"
)

/**
只能加不能减少
*/
var keyMap = []string{
	"sgMExvjDEzpsbei5",
	"zIdOBazRNmFpZxce",
	"0nTQXY6GVGnnEtqR",
	"OL1I4A7JxDQMEk6V",
	"JbBhh1AR9W894jJx",
	"QKnf9YgBpJD1hdN3",
	"SFQdtzWsKdG7IwPh",
	"uW71BbGD7FsKaMOK",
	"VNHbC2dStDZ6ofm5",
	"jUe8BLykf1FvCEoG",
	"XzSkNc9tUjN372eV",
	"uS2WcCNlm0pglTCK",
	"1Mjz9w1T4AhrlnXO",
	"gC83cwcaIR8AxwJE",
	"1h20YAS1b10ctJtn",
	"jNzTtI16HJhWgqo7",
}

type InvalidKeyIndexError string

func (i InvalidKeyIndexError) Error() string {
	return "invalid key index: " + string(i)
}

func getKey(keyIndex int) (string, error) {
	if keyIndex >= 0 && keyIndex < len(keyMap) {
		return keyMap[keyIndex], nil
	}
	return "", InvalidKeyIndexError(keyIndex)
}

func AesEncrypt(orig string) (string, int, error) {
	keyIndex := rand.Intn(16)
	key, err := getKey(keyIndex)
	if err != nil {
		return "", keyIndex, err
	}
	// 转成字节数组
	origData := []byte(orig)
	k := []byte(key)

	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)

	return base64.StdEncoding.EncodeToString(cryted), keyIndex, nil

}

func AesDecrypt(cryted string, keyIndex int) (string, error) {
	key, err := getKey(keyIndex)
	if err != nil {
		return "", err
	}
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
	k := []byte(key)

	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = PKCS7UnPadding(orig)
	return string(orig), nil
}

//补码
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//去码
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
