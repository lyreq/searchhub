package slugger

import (
	"regexp"
	"strings"
)

func Slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "-")

	// Türkçe karakterleri temizle
	replacer := strings.NewReplacer(
		"ç", "c", "ğ", "g", "ı", "i", "ö", "o", "ş", "s", "ü", "u",
	)
	s = replacer.Replace(s)

	// Alfasayısal ve tire dışındaki karakterleri sil
	re := regexp.MustCompile(`[^a-z0-9\-]`)
	return re.ReplaceAllString(s, "")
}
