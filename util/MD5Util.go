package util

import (
	"crypto/md5"
	"fmt"
)

func MD5Encrypt(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has)
	return md5str1
}
