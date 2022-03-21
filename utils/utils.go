package utils

import (
	"crypto/rand"
	"math/big"
)

var Letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

// 生成随机字符串，长度为n
func GenRandString(n int) string {
	arr := make([]rune, n)
	for i := range arr {
		lenght := int64(len(Letters))
		randint, _ := rand.Int(rand.Reader, big.NewInt(lenght))
		arr[i] = Letters[randint.Int64()]
	}
	return string(arr)
}
