package utils

import (
	"math/rand"
	"os/exec"
	"time"
)

// GetRandomString 用于生成指定长度的随机字符串.
func GetRandomString(len int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < len; i++ {
		result = append(result, bytes[r.Intn(62)])
	}
	return string(result)
}

// GetUUID 获取UUID(唯一性质的字符串标识).
func GetUUID() (string, error) {
	uuid, err := exec.Command("uuidgen").Output()
	if err != nil {
		return "", err
	}
	return Get32MD5Encode(string(uuid)), nil
}
