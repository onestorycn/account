package library

import (
	"math/rand"
	"time"
	"bytes"
	"crypto/md5"
	"encoding/hex"
)

func RandomString(l int) string {
	var result bytes.Buffer
	var temp string
	for i := 0; i < l; {
		if string(RandInt(65, 90)) != temp {
			temp = string(RandInt(65, 90))
			result.WriteString(temp)
			i++
		}
	}
	return result.String()
}
func RandInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func IsEmpty(input interface{}) bool {
	switch input.(type) {
	case int:
		if input == 0 {
			return true
		}
		break
	case int32:
		if input == 0 {
			return true
		}
		break
	case int64:
		if input == 0 {
			return true
		}
		break
	case string:
		if input == "" || len([]rune(input.(string))) < 1 {
			return true
		}
		break
	default:
		return false
		break
	}
	return false
}

func EncodeMd5(rawVal string) string {
	h := md5.New()
	h.Write([]byte(rawVal))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}