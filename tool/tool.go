package tool

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"

	uuid "github.com/nu7hatch/gouuid"
	"google.golang.org/grpc/status"
)

// CreateSha1 generate sha1 from source
func CreateSha1(source string) string {
	data := []byte(source)
	return fmt.Sprintf("%x", sha1.Sum(data))
}

// CreateMD5 generate sha1 from source
func CreateMD5(source string) string {
	data := []byte(source)
	return fmt.Sprintf("%x", md5.Sum(data))
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
	val, err := strconv.ParseFloat(source, 64)
	if err != nil {
		return def
	}
	return val
}

// PadRight ("go", "x", 6) => "goxxxx"
func PadRight(str, pad string, lenght int) string {
	for {
		str += pad
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}

// PadLeft ("go", "x", 6) => "xxxxgo"
func PadLeft(str, pad string, lenght int) string {
	for {
		str = pad + str
		if len(str) >= lenght {
			return str[0:lenght]
		}
	}
}

// PrintError from converted status then mark test as fail and then exit
func PrintError(t *testing.T, title string, err error) {
	s := status.Convert(err)
	fmt.Printf("%s: %s\n\n", title, s.Message())
	t.FailNow()
}

// PrintStruct helper
func PrintStruct(title, data interface{}) {
	fmt.Printf("%s:\n%s\n\n", title, AsJSON(data))
}

// StringToNumberString fix string to
func StringToNumberString(val string) string {
	v, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return "0"
	}
	return strconv.FormatInt(v, 10)
}

// SumString number sum of two string
func SumString(a, b string) string {
	aa, _ := strconv.ParseInt(a, 10, 64)
	bb, _ := strconv.ParseInt(b, 10, 64)
	return strconv.FormatInt(aa+bb, 10)
}

// WordWraps long string into slice with maximum char length
func WordWraps(source string, maxWidth int) []string {
	var result = make([]string, 0)
	var wrapped string

	// Split string into array of words
	words := strings.Fields(source)

	if len(words) == 0 {
		return result
	}

	remaining := maxWidth

	for _, word := range words {
		if len(word)+1 > remaining {
			if len(wrapped) > 0 {
				result = append(result, wrapped)
				wrapped = ""
			}

			wrapped += word
			remaining = maxWidth - len(word)
		} else {
			if len(wrapped) > 0 {
				wrapped += " "
			}

			wrapped += word
			remaining = remaining - (len(word) + 1)
		}
	}

	if len(wrapped) > 0 {
		result = append(result, wrapped)
	}

	return result
}
