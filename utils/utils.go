package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math"
	"math/rand"
)

// MD5String The returned value is a 32-bit lowercase hexadecimal string
func MD5String(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

// Return amount*10^6
func ToUnit(amount float64) float64 {
	return amount * math.Pow10(6)
}

// RandInt 生成一个6位随机数
func RandInt(max int) int {
	return rand.Intn(max)
}