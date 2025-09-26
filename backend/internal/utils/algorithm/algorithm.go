package algorithm

import (
	"fmt"
	"math/rand"
	"regexp"
	"runtime"

	"golang.org/x/crypto/bcrypt"
)

func GenerateApiUsername() string {
	return RandStringBytes(36, SmallMixedBytes)
}

func GenerateApiPassword() string {
	return RandStringBytes(36, SmallMixedBytes)
}

func RandNumber(start int, end int) int {
	code := rand.Intn(end-start) + start
	return code
}

type Charset string

const (
	SmallMixedBytes  Charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	MediumMixedBytes Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

func RandStringBytes(n int, charset Charset) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func GetTrace(skip int) (line string, functionName string) {
	pc, file, lineNumber, ok := runtime.Caller(skip)
	if !ok {
		return "?", "?"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return fmt.Sprintf("%s:%d", file, lineNumber), "?"
	}

	return fmt.Sprintf("%s:%d", file, lineNumber), fn.Name()
}

func GetFunctionName(skip int) string {
	_, functionName := GetTrace(skip)
	return functionName
}

func GetLineTrace(skip int) string {
	line, _ := GetTrace(skip)
	return line
}

func IsPasswordCorrect(hash, pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)) == nil
}

func GeneratePassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashedPassword), err
}
func IsValidPassword(password string) bool {
	if len(password) < 6 {
		return false
	}

	uppercase := regexp.MustCompile(`[A-Z]`)
	lowercase := regexp.MustCompile(`[a-z]`)
	special := regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]{};':"\\|,.<>\/?]`)

	return uppercase.MatchString(password) &&
		lowercase.MatchString(password) &&
		special.MatchString(password)
}

func SliceSafeGet[T any](slice []T, index int, fallback T) T {
	if index < 0 || index >= len(slice) {
		return fallback
	}
	return slice[index]
}

func Find[T any](slice []T, predicate func(T) bool) (T, bool) {
	for _, item := range slice {
		if predicate(item) {
			return item, true
		}
	}
	return *new(T), false
}
