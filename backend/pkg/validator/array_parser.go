// validator/array_parser.go (yeni)
package validator

import (
	"encoding/json"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func (v *validatorImpl) formArrayParser(c echo.Context, o interface{}) error {
	// 1) Önce multipart mı diye dene (dosyalar burada)
	var values url.Values
	if mf, err := c.MultipartForm(); err == nil && mf != nil {
		values = mf.Value
	} else {
		// 2) Değilse normal form'u parse et
		_ = c.Request().ParseForm()
		values = c.Request().Form
	}

	rv := reflect.ValueOf(o)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return nil
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return nil
	}
	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		tag := f.Tag.Get("form")
		if tag == "" || tag == "-" {
			continue
		}
		fv := rv.Field(i)
		if !fv.CanSet() || fv.Kind() != reflect.Slice {
			continue
		}

		elemKind := fv.Type().Elem().Kind()
		if !(elemKind == reflect.String ||
			elemKind == reflect.Int || elemKind == reflect.Int64 || elemKind == reflect.Int32 || elemKind == reflect.Int16 || elemKind == reflect.Int8 ||
			elemKind == reflect.Uint || elemKind == reflect.Uint64 || elemKind == reflect.Uint32 || elemKind == reflect.Uint16 || elemKind == reflect.Uint8) {
			continue
		}

		// 1) Aynı key tekrarları
		candidates := []string{}
		if arr, ok := values[tag]; ok && len(arr) > 0 {
			candidates = append(candidates, arr...)
		}
		// 2) Laravel-vari key: tag+"[]"
		if arr, ok := values[tag+"[]"]; ok && len(arr) > 0 {
			candidates = append(candidates, arr...)
		}
		// 3) Tek satır CSV veya JSON için tek değer yakala
		if len(candidates) == 0 {
			one := values.Get(tag)
			if one == "" {
				one = values.Get(tag + "[]")
			}
			if one != "" {
				candidates = append(candidates, one)
			}
		}

		if len(candidates) == 0 {
			continue
		}

		// Tek eleman JSON array ise: [1,2] / ["a","b"]
		if len(candidates) == 1 && strings.HasPrefix(strings.TrimSpace(candidates[0]), "[") {
			raw := candidates[0]
			switch elemKind {
			case reflect.String:
				var tmp []string
				if json.Unmarshal([]byte(raw), &tmp) == nil {
					fv.Set(reflect.ValueOf(tmp))
				}
			default:
				// sayısal
				var ints []int64
				if json.Unmarshal([]byte(raw), &ints) == nil {
					out := reflect.MakeSlice(fv.Type(), 0, len(ints))
					for _, n := range ints {
						out = reflect.Append(out, reflect.ValueOf(n).Convert(fv.Type().Elem()))
					}
					fv.Set(out)
				} else {
					var uints []uint64
					if json.Unmarshal([]byte(raw), &uints) == nil {
						out := reflect.MakeSlice(fv.Type(), 0, len(uints))
						for _, n := range uints {
							out = reflect.Append(out, reflect.ValueOf(n).Convert(fv.Type().Elem()))
						}
						fv.Set(out)
					} else {
						var floats []float64
						if json.Unmarshal([]byte(raw), &floats) == nil {
							out := reflect.MakeSlice(fv.Type(), 0, len(floats))
							for _, n := range floats {
								out = reflect.Append(out, castNumber(elemKind, n))
							}
							fv.Set(out)
						}
					}
				}
			}
			continue
		}

		// CSV desteği: "15,27,42"
		expanded := make([]string, 0, len(candidates))
		for _, cval := range candidates {
			parts := strings.Split(cval, ",")
			for _, p := range parts {
				if s := strings.TrimSpace(p); s != "" {
					expanded = append(expanded, s)
				}
			}
		}

		out := reflect.MakeSlice(fv.Type(), 0, len(expanded))
		for _, s := range expanded {
			switch elemKind {
			case reflect.String:
				out = reflect.Append(out, reflect.ValueOf(s))
			default:
				if strings.HasPrefix(s, "+") {
					s = strings.TrimPrefix(s, "+")
				}
				if strings.HasPrefix(s, "-") && (elemKind == reflect.Uint || elemKind == reflect.Uint64 || elemKind == reflect.Uint32 || elemKind == reflect.Uint16 || elemKind == reflect.Uint8) {
					continue
				}
				if strings.ContainsAny(s, ".eE") {
					continue
				}
				if elemKind == reflect.Int || elemKind == reflect.Int64 || elemKind == reflect.Int32 || elemKind == reflect.Int16 || elemKind == reflect.Int8 {
					if i, err := strconv.ParseInt(s, 10, 64); err == nil {
						out = reflect.Append(out, reflect.ValueOf(i).Convert(fv.Type().Elem()))
					}
				} else {
					if u, err := strconv.ParseUint(s, 10, 64); err == nil {
						out = reflect.Append(out, reflect.ValueOf(u).Convert(fv.Type().Elem()))
					}
				}
			}
		}
		fv.Set(out)
	}
	return nil
}

func castNumber(kind reflect.Kind, n float64) reflect.Value {
	switch kind {
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return reflect.ValueOf(int64(n))
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		if n < 0 {
			n = 0
		}
		return reflect.ValueOf(uint64(n))
	default:
		return reflect.Value{}
	}
}
