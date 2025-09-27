package utils

import (
	"errors"
	"fmt"
	"lytemp/internal/domain/dtos/responses"
	"lytemp/internal/domain/models"
	"lytemp/pkg/provider"
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

func StrconvI(v int, def int) string {
	if v <= 0 {
		v = def
	}
	return strconv.Itoa(v)
}


func TypeCoefFor(t string) float64 {
	switch t {
	case string(models.ContentTypeVideo):
		return 1.5
	case string(models.ContentTypeArticle):
		return 1.0
	default:
		return 1.0
	}
}

func DefInt(v, d int) int {
	if v <= 0 {
		return d
	}
	return v
}

func BaseScore(it provider.Item) float64 {
	switch StringsLower(it.NormalizedType) {
	case "video":
		return float64(Deref(it.Views))/1000.0 + float64(Deref(it.Likes))/100.0
	case "article":
		return float64(Deref(it.ReadingTimeMinutes)) + float64(Deref(it.Reactions))/50.0
	default:
		return 0
	}
}

func TypeCoef(t string) float64 {
	switch StringsLower(t) {
	case "video":
		return 1.5
	case "article":
		return 1.0
	default:
		return 1.0
	}
}

func FreshnessScore(pub time.Time, now time.Time) float64 {
	days := now.Sub(pub).Hours() / 24.0
	switch {
	case days <= 7:
		return 5
	case days <= 30:
		return 3
	case days <= 90:
		return 1
	default:
		return 0
	}
}

func EngagementScore(it provider.Item) float64 {
	switch StringsLower(it.NormalizedType) {
	case "video":
		views := float64(Max1(Deref(it.Views)))
		likes := float64(Deref(it.Likes))
		return (likes / views) * 10.0
	case "article":
		rt := float64(Max1(Deref(it.ReadingTimeMinutes)))
		re := float64(Deref(it.Reactions))
		return (re / rt) * 5.0
	default:
		return 0
	}
}

func Deref(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}
func Max1(x int) int {
	if x < 1 {
		return 1
	}
	return x
}
func StringsLower(s string) string {
	// tiny local helper to avoid importing strings in each function
	for _, b := range s {
		if b >= 'A' && b <= 'Z' {
			// shadowed to simple lower; for ASCII it's fine
		}
	}
	// for correctness import strings:
	// return strings.ToLower(s)
	// But we actually need real lower; let's just import strings.
	return ToLower(s)
}
func ToLower(s string) string { // real impl via strings
	return strings.ToLower(s)
}
func Round2(f float64) float64 {
	return float64(int64(f*100+0.5)) / 100.0
}
func DetectFormat(name string) models.ProviderFormat {
	if name == "xml_provider" {
		return models.ProviderFormatXML
	}
	return models.ProviderFormatJSON
}
