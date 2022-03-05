package base64

import "encoding/base64"

// Encoding 将指定字节编码为 Base64 字符串格式.
func Encoding(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// Decoding 将一个 Base64 格式的字符串进行解码为字节类型.
func Decoding(base64Str string) []byte {
	bytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return []byte("")
	}
	return bytes
}
