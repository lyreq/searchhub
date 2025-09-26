package date

import (
	"fmt"
	"strings"
	"time"
)

// Türkçe ay isimleri
var TurkishMonths = map[int]string{
	1:  "Ocak",
	2:  "Şubat",
	3:  "Mart",
	4:  "Nisan",
	5:  "Mayıs",
	6:  "Haziran",
	7:  "Temmuz",
	8:  "Ağustos",
	9:  "Eylül",
	10: "Ekim",
	11: "Kasım",
	12: "Aralık",
}

// Türkçe gün isimleri
var TurkishDays = map[time.Weekday]string{
	time.Sunday:    "Pazar",
	time.Monday:    "Pazartesi",
	time.Tuesday:   "Salı",
	time.Wednesday: "Çarşamba",
	time.Thursday:  "Perşembe",
	time.Friday:    "Cuma",
	time.Saturday:  "Cumartesi",
}

// FormatTurkishDate Türkçe tarih formatı oluşturur
func FormatTurkishDate(date time.Time) string {
	day := date.Day()
	monthName := TurkishMonths[int(date.Month())]
	dayName := TurkishDays[date.Weekday()]
	return strings.Join([]string{fmt.Sprintf("%d", day), monthName, dayName}, " ")
}

// ParseDate tarih string'ini parse eder
func ParseDate(dateStr string) (time.Time, error) {
	if dateStr == "bugün" {
		return time.Now(), nil
	}
	return time.Parse("02/01/2006", dateStr)
}

// RetypeMessage mesaj şablonundaki placeholder'ları gerçek değerlerle değiştirir
func RetypeMessage(message string, typeInfo map[string]string) string {
	message = strings.ReplaceAll(message, "[tarih]", typeInfo["date"])
	message = strings.ReplaceAll(message, "[ders adı]", typeInfo["lesson"])
	message = strings.ReplaceAll(message, "[gec-gelmedi]", typeInfo["status"])
	message = strings.ReplaceAll(message, "[ders saati]", typeInfo["hour"])
	return message
}
