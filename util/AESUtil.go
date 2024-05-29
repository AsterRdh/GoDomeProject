package util

import (
	"awesomeProject/model"
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"errors"
)

// AesEncryptByECB 加密
func AesEncryptByECB(data, key string) string {
	// 判断key长度
	keyLenMap := map[int]struct{}{16: {}, 24: {}, 32: {}}
	if _, ok := keyLenMap[len(key)]; !ok {
		panic("key长度必须是 16、24、32 其中一个")
	}
	// 密钥和待加密数据转成[]byte
	originByte := []byte(data)
	keyByte := []byte(key)
	// 创建密码组，长度只能是16、24、32字节
	block, _ := aes.NewCipher(keyByte)
	// 获取密钥长度
	blockSize := block.BlockSize()
	// 补码
	originByte = PKCS7Padding(originByte, blockSize)
	// 创建保存加密变量
	encryptResult := make([]byte, len(originByte))
	// CEB是把整个明文分成若干段相同的小段，然后对每一小段进行加密
	for bs, be := 0, blockSize; bs < len(originByte); bs, be = bs+blockSize, be+blockSize {
		block.Encrypt(encryptResult[bs:be], originByte[bs:be])
	}
	return base64.StdEncoding.EncodeToString(encryptResult)
}

// PKCS7Padding 补码
func PKCS7Padding(originByte []byte, blockSize int) []byte {
	// 计算补码长度
	padding := blockSize - len(originByte)%blockSize
	// 生成补码
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	// 追加补码
	return append(originByte, padText...)
}

// 解密
func AesDecryptByECB(data, key string) (string, string, error) {
	// 判断key长度
	keyLenMap := map[int]struct{}{16: {}, 24: {}, 32: {}}
	if _, ok := keyLenMap[len(key)]; !ok {
		return "", "key长度必须是 16、24、32 其中一个", errors.New("key长度异常")
	}
	// 反解密码base64
	originByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", "反解密码base64失败", err
	}

	// 密钥和待加密数据转成[]byte
	keyByte := []byte(key)
	// 创建密码组，长度只能是16、24、32字节
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", "创建密码组创建失败", err
	}
	// 获取密钥长度
	blockSize := block.BlockSize()
	// 创建保存解密变量
	decrypted := make([]byte, len(originByte))
	for bs, be := 0, blockSize; bs < len(originByte); bs, be = bs+blockSize, be+blockSize {
		block.Decrypt(decrypted[bs:be], originByte[bs:be])
	}
	// 解码
	return string(PKCS7UNPadding(decrypted)), "", nil
}

// PKCS7UNPadding 解码
func PKCS7UNPadding(originDataByte []byte) []byte {
	length := len(originDataByte)
	padding := int(originDataByte[length-1])
	return originDataByte[:(length - padding)]
}

func GetAesKey(salt string) string {
	var AESKey = MD5Encrypt(salt + model.AESKey)
	return AESKey
}
