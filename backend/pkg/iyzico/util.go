package iyzico

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
)

type Charset string

const (
	UPPER_ALPHABET Charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LOWER_ALPHABET Charset = "abcdefghijklmnopqrstuvwxyz"
	NUMERIC        Charset = "0123456789"

	ALPHANUMERIC Charset = UPPER_ALPHABET + NUMERIC
)

func RandomString(charset Charset, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = byte(charset[rand.Intn(len(charset))])
	}
	return string(b)
}

func GenerateRequestString(body interface{}) string {
	if false {
		a, _ := json.Marshal(body)
		return string(a)
	}

	val := reflect.ValueOf(body)
	typ := reflect.TypeOf(body)

	str := "["
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldTyp := typ.Field(i)
		tag := fieldTyp.Tag.Get("json")

		str += fmt.Sprintf("%s=", tag)

		switch fieldVal.Kind() {
		case reflect.Slice:
			str += "["
			for j := 0; j < fieldVal.Len(); j++ {
				str += GenerateRequestString(fieldVal.Index(j).Interface())
				str += ", "
			}
			if fieldVal.Len() > 0 {
				str = str[:len(str)-2]
			}
			str += "]"
		case reflect.Struct:
			str += GenerateRequestString(fieldVal.Interface())
		default:
			str += fmt.Sprintf("%v", fieldVal.Interface())
		}

		str += ","
	}

	if len(str) > 0 {
		str = str[:len(str)-1]
	}
	str += "]"

	return str
}

func GenerateHash(str string) string {
	return base64.StdEncoding.EncodeToString(hashSHA1(str))
}

func hashSHA1(str string) []byte {
	hash := sha1.New()
	hash.Write([]byte(str))
	return hash.Sum(nil)
}

func GetHTTPHeaders(str string, randStr string) map[string]string {
	m := make(map[string]string)

	m["Content-Type"] = "application/json"
	m["x-iyzi-rnd"] = randStr
	m["Authorization"] = str

	return m
}
