package phone

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/nyaruka/phonenumbers"
)

// (2|3|5|8) + 9 adet rakam
//
// Kabul edilebilir numaralar: 2XXXXXXXXX, 3XXXXXXXXX, 5XXXXXXXXX, 8XXXXXXXXX
func IsPhoneNumberValid(phone string) bool {
	re := regexp.MustCompile(`^[2358]\d{9}$`)
	return re.MatchString(phone)
}

type PhoneNumber struct {
	NationalPhone string
	CountryCode   string
}

func GetCountryCode(phone string) (string, error) {
	num, err := phonenumbers.Parse(phone, "")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("+%d", num.GetCountryCode()), nil
}

func GetFullPhoneNumber(phone, countryCode string) string {
	return fmt.Sprintf("+%s%s", countryCode, phone)
}

func IsValidPhoneNumber(phone string) bool {
	num, err := phonenumbers.Parse(phone, "")
	if err != nil {
		return false
	}
	return phonenumbers.IsValidNumber(num)
}

func GetNationalNumber(phone string) (string, error) {
	num, err := phonenumbers.Parse(phone, "")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", num.GetNationalNumber()), nil
}

func GetIsoCountryCodeFromPhone(phone string) (string, error) {
	num, err := phonenumbers.Parse(phone, "")
	if err != nil {
		return "", err
	}
	return phonenumbers.GetRegionCodeForNumber(num), nil
}

func ParseNumber(phone string) (PhoneNumber, error) {
	num, err := phonenumbers.Parse(phone, "")
	if err != nil {
		return PhoneNumber{}, err
	}

	return PhoneNumber{
		NationalPhone: strconv.Itoa(int(num.GetNationalNumber())),
		CountryCode:   strconv.Itoa(int(num.GetCountryCode())),
	}, nil
}
func ValidatePhoneWithDialCode(phone string) bool {
	// dial code +90 ve 10 haneli telefon numarası
	re := regexp.MustCompile(`^\+90\d{10}$`)
	return re.MatchString(phone)
}

func ValidateAndNormalizeTRPhone(phone string) (string, error) {
	num, err := phonenumbers.Parse(phone, "TR")
	if err != nil {
		return "", errors.New("Telefon numarası geçersiz.")
	}

	if !phonenumbers.IsValidNumber(num) {
		return "", errors.New("Telefon numarası geçerli değil.")
	}

	if phonenumbers.GetRegionCodeForNumber(num) != "TR" {
		return "", errors.New("Sadece Türkiye numaralarına izin verilmektedir.")
	}

	return fmt.Sprintf("%d", num.GetNationalNumber()), nil
}

func ParseTRNumber(phone string) (PhoneNumber, error) {
	num, err := phonenumbers.Parse(phone, "TR")
	if err != nil {
		return PhoneNumber{}, err
	}

	if !phonenumbers.IsValidNumber(num) {
		return PhoneNumber{}, errors.New("telefon numarası geçerli değil")
	}

	return PhoneNumber{
		NationalPhone: strconv.Itoa(int(num.GetNationalNumber())),
		CountryCode:   strconv.Itoa(int(num.GetCountryCode())),
	}, nil
}
