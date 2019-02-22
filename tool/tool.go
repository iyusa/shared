package tool

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strconv"

	uuid "github.com/nu7hatch/gouuid"
)

// CreateSha1 generate sha1 from source
func CreateSha1(source string) string {
	data := []byte(source)
	return fmt.Sprintf("%x", sha1.Sum(data))
}

// CreateStan create random string uuid, example: 66889689-3a29-4104-6bbb-e13782d36b1d
func CreateStan() string {
	u4, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	return u4.String()
}

// AsJSON convert struct to json string
func AsJSON(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Sprintf("Gagal konversi: %s", err.Error())
	}
	return string(b)
}

// Converter

// StrToInt convert string to integer, return def if fail
func StrToInt(source string, def int) int {
	val, err := strconv.Atoi(source)
	if err != nil {
		return def
	}
	return val
}

// StrToFloat convert string to float, return def if fail
func StrToFloat(source string, def float64) float64 {
	val, err := strconv.ParseFloat(source, 32)
	if err != nil {
		return def
	}
	return val
}
