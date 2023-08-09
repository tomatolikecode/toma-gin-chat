package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// 小写md5
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))

	return hex.EncodeToString(h.Sum(nil))
}

// 大写md5
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

// 密码加密
func MakePassword(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)
}

// 密码解密
func ValidPassword(plainpwd, salt string, password string) bool {
	return Md5Encode(plainpwd+salt) == password
}
