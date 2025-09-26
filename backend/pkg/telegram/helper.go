package telegram

import (
	"strings"
)

// Helper represents a Telegram message helper
type Helper struct {
	message string
}

// NewHelper creates a new Telegram helper instance
func NewHelper() *Helper {
	return &Helper{}
}

func (h *Helper) GenerateElockRegisterMessage(data map[string]interface{}) string {
	h.SetTitle("lytemp E-KİLİT KAYIT BİLDİRİMİ").
		SetOption("Kurum Kodu", data["company_code"].(string), true).
		SetOption("Okul/Kurum Adı", data["company_name"].(string), true).
		SetMessage("—————————\n").
		SetOption("lytemp Üye ID", data["id"].(string), true).
		SetOption("Adı", data["name"].(string), true).
		SetOption("Telefon Numarası", data["phone"].(string), true).
		SetOption("İl/İlçe", data["city"].(string)+"/"+data["district"].(string), true)

	return h.GetMessage()
}
func (h *Helper) GenerateInformationRequestMessage(data map[string]interface{}) string {
	h.SetTitle("- Bilgi Talebi - lytemp.com").
		SetOption("Kurum Kodu", data["company_code"].(string), true).
		SetOption("Okul/Kurum Adı", data["company_name"].(string), true).
		SetMessage("—————————\n").
		SetOption("lytemp Üye ID", data["id"].(string), true).
		SetOption("Adı", data["name"].(string), true).
		SetOption("Telefon Numarası", data["phone"].(string), true).
		SetOption("İlgilenilen Projeler", data["project"].(string), true)

	return h.GetMessage()
}
func (h *Helper) GeneratelytempToMySchoolMessage(data map[string]interface{}) string {
	h.SetTitle("Okulum - İlksms'ten geçiş yapan kullanıcı").
		SetOption("Kullanıcı Adı", data["username"].(string), true).
		SetOption("Okul/Kurum Adı", data["company_name"].(string), true).
		SetMessage("—————————\n").
		SetOption("lytemp Üye ID", data["id"].(string), true).
		SetOption("Eposta", data["email"].(string), true).
		SetOption("Telefon Numarası", data["phone"].(string), true)
	return h.GetMessage()
}

// SetTitle sets the title of the message
func (h *Helper) SetTitle(title string) *Helper {
	h.SetMessage("*- " + title + " -*\n\n")
	return h
}

// SetOption adds an option to the message
func (h *Helper) SetOption(option, message string, clean bool) *Helper {
	if clean {
		message = strings.TrimSpace(message)
	}
	h.SetMessage("*" + option + ":* " + message + "\n")
	return h
}

// SetMessage appends to the current message
func (h *Helper) SetMessage(m string) *Helper {
	h.message += m
	return h
}

// GetMessage returns the current message
func (h *Helper) GetMessage() string {
	return h.message
}
