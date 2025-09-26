package utils

import (
	"lytemp/internal/domain/dtos/responses"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func BasicError(err interface{}, message ...string) responses.Response[string] {
	msg := ""
	switch err.(type) {
	case error:
		if len(message) > 0 {
			msg = message[0]
		} else {
			msg = "Beklenmeyen bir hata oluştu lütfen canlı destek ile iletişime geçiniz. Hata: " + err.(error).Error()
		}

		if errors.Is(err.(error), gorm.ErrRecordNotFound) {

			if len(message) > 0 {
				msg = message[0]
			} else {
				msg = "Veritabanında kayıt bulunamadı."
			}
		}

	case string:
		msg = err.(string)
	case nil:
		msg = message[0]
	default:
		msg = "Bilinmeyen bir hata oluştu."
	}

	return responses.Error("", msg)
}

// Validate validates the request body
//
// Deprecated: Use validator.Validator instead
func Validate(c *echo.Context, req interface{}) error {
	_ = (*c).Bind(req)
	err := (*c).Validate(req)
	return err
}

func DerefStr(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// StringPtrOrNil returns *string only if s is not empty, otherwise returns nil
func StringPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func IntPtrToString(p *int) string {
	if p == nil {
		return ""
	}
	return strconv.Itoa(*p)
}

func Int32PtrToString(p *int32) string {
	if p == nil {
		return ""
	}
	return strconv.FormatInt(int64(*p), 10)
}

func Int64PtrToString(p *int64) string {
	if p == nil {
		return ""
	}
	return strconv.FormatInt(*p, 10)
}

func ParseOptionalTime(val string, loc *time.Location, layouts ...string) (*time.Time, error) {
	val = strings.TrimSpace(strings.Trim(val, `"`)) // tırnak veya boşluk temizle
	if val == "" {
		return nil, nil
	}
	if loc == nil {
		loc = time.UTC
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, val, loc); err == nil {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("desteklenmeyen tarih formatı: %q", val)
}
