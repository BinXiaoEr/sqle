package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"unsafe"
)

// base64 encoding string to decode string

func DecodeString(base64Str string) (string, error) {
	sDec, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return "", err
	}
	return *(*string)(unsafe.Pointer(&sDec)), nil
}

func Md5String(data string) string {
	md5 := md5.New()
	md5.Write([]byte(data))
	md5Data := md5.Sum([]byte(nil))
	return hex.EncodeToString(md5Data)
}

func HasPrefix(s, prefix string, caseSensitive bool) bool {
	if caseSensitive {
		return strings.HasPrefix(s, prefix)
	}
	return strings.HasPrefix(strings.ToLower(s), strings.ToLower(prefix))
}

func HasSuffix(s, suffix string, caseSensitive bool) bool {
	if caseSensitive {
		return strings.HasSuffix(s, suffix)
	}
	return strings.HasSuffix(strings.ToLower(s), strings.ToLower(suffix))
}

func GetDuplicate(c []string) []string {
	d := []string{}
	for i, v1 := range c {
		for j, v2 := range c {
			if i >= j {
				continue
			}
			if v1 == v2 {
				d = append(d, v1)
			}
		}
	}
	return RemoveDuplicate(d)
}

func RemoveDuplicate(c []string) []string {
	var tmpMap = map[string]struct{}{}
	var result = []string{}
	for _, v := range c {
		beforeLen := len(tmpMap)
		tmpMap[v] = struct{}{}
		AfterLen := len(tmpMap)
		if beforeLen != AfterLen {
			result = append(result, v)
		}
	}
	return result
}
