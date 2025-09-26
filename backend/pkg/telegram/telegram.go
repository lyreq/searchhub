package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"lytemp/config"
	"net/http"
	"net/url"
)

// ChatKey represents available chat keys
type ChatKey struct {
	ErrorLogger                    string
	SmsMemberships                 string
	ErrorNotification              string
	GeneralNotifications           string
	Turnstile                      string
	SpeedTest                      string
	BugErrors                      string
	PackageNotification            string
	SmsError                       string
	EschoolAppointmentNotification string
	TechnicalSupport               string
}

// Chat represents a Telegram chat configuration
type Chat struct {
	ID string
}

// Telegram represents the Telegram API client
type Telegram struct {
	chats map[string]Chat
	token string
	url   string
	Keys  ChatKey
}

// New creates a new Telegram instance
func New() (*Telegram, error) {
	t := &Telegram{
		chats: map[string]Chat{
			"error-logger":           {ID: "-701969369"},
			"sms-uyelikleri":         {ID: "-1001651719184"},
			"hata-bildirimi":         {ID: "-610417756"},
			"genel-bildirimler":      {ID: "-1001701146778"},
			"turnike":                {ID: "-1001605819520"},
			"hiz-testi":              {ID: "-4171864813"},
			"bug-errors":             {ID: "-100824491195"},
			"paket-bildirim":         {ID: "-1001885490087"},
			"sms-error":              {ID: "-4237070771"},
			"eokul-randevu-bildirim": {ID: "-4067639239"},
			"teknik-destek":          {ID: "-1001646781453"},
		},
		url: "https://api.telegram.org/bot",
		Keys: ChatKey{
			ErrorLogger:                    "error-logger",
			SmsMemberships:                 "sms-uyelikleri",
			ErrorNotification:              "hata-bildirimi",
			GeneralNotifications:           "genel-bildirimler",
			Turnstile:                      "turnike",
			SpeedTest:                      "hiz-testi",
			BugErrors:                      "bug-errors",
			PackageNotification:            "paket-bildirim",
			SmsError:                       "sms-error",
			EschoolAppointmentNotification: "eokul-randevu-bildirim",
			TechnicalSupport:               "teknik-destek",
		},
	}

	if err := t.setToken(); err != nil {
		return nil, err
	}

	return t, nil
}

// Message sends a text message to the specified chat
func (t *Telegram) Message(chatKey string, message string, params map[string]interface{}) (map[string]interface{}, error) {
	if message == "" {
		return nil, fmt.Errorf("message not valid")
	}

	chat, exists := t.chats[chatKey]
	if !exists {
		return nil, fmt.Errorf("telegram chat not found: %s", chatKey)
	}

	parameters := url.Values{}
	parameters.Set("chat_id", chat.ID)
	parameters.Set("text", message)
	parameters.Set("parse_mode", "Markdown")
	for key, value := range params {
		parameters.Set(key, fmt.Sprintf("%v", value))
	}

	return t.Send("sendMessage", parameters)
}

// MessageTo sends a message to a specific chat using ChatKey
func (t *Telegram) MessageTo(key string, message string, params map[string]interface{}) (map[string]interface{}, error) {
	var chatKey string
	switch key {
	case t.Keys.ErrorLogger:
		chatKey = t.Keys.ErrorLogger
	case t.Keys.SmsMemberships:
		chatKey = t.Keys.SmsMemberships
	case t.Keys.ErrorNotification:
		chatKey = t.Keys.ErrorNotification
	case t.Keys.GeneralNotifications:
		chatKey = t.Keys.GeneralNotifications
	case t.Keys.Turnstile:
		chatKey = t.Keys.Turnstile
	case t.Keys.SpeedTest:
		chatKey = t.Keys.SpeedTest
	case t.Keys.BugErrors:
		chatKey = t.Keys.BugErrors
	case t.Keys.PackageNotification:
		chatKey = t.Keys.PackageNotification
	case t.Keys.SmsError:
		chatKey = t.Keys.SmsError
	case t.Keys.EschoolAppointmentNotification:
		chatKey = t.Keys.EschoolAppointmentNotification
	case t.Keys.TechnicalSupport:
		chatKey = t.Keys.TechnicalSupport
	default:
		return nil, fmt.Errorf("invalid chat key: %s", key)
	}

	return t.Message(chatKey, message, params)
}

// Contact sends contact information to the specified chat
func (t *Telegram) Contact(chatKey string, phone, firstName, lastName string) (map[string]interface{}, error) {
	if phone == "" || firstName == "" || lastName == "" {
		return nil, fmt.Errorf("contact parameters not valid")
	}

	chat, exists := t.chats[chatKey]
	if !exists {
		return nil, fmt.Errorf("telegram chat not found: %s", chatKey)
	}

	parameters := url.Values{}
	parameters.Set("chat_id", chat.ID)
	parameters.Set("phone_number", phone)
	parameters.Set("first_name", firstName)
	parameters.Set("last_name", lastName)

	return t.Send("sendContact", parameters)
}

// ContactTo sends contact information to a specific chat using ChatKey
func (t *Telegram) ContactTo(key string, phone, firstName, lastName string) (map[string]interface{}, error) {
	var chatKey string
	switch key {
	case t.Keys.ErrorLogger:
		chatKey = t.Keys.ErrorLogger
	case t.Keys.SmsMemberships:
		chatKey = t.Keys.SmsMemberships
	case t.Keys.ErrorNotification:
		chatKey = t.Keys.ErrorNotification
	case t.Keys.GeneralNotifications:
		chatKey = t.Keys.GeneralNotifications
	case t.Keys.Turnstile:
		chatKey = t.Keys.Turnstile
	case t.Keys.SpeedTest:
		chatKey = t.Keys.SpeedTest
	case t.Keys.BugErrors:
		chatKey = t.Keys.BugErrors
	case t.Keys.PackageNotification:
		chatKey = t.Keys.PackageNotification
	case t.Keys.SmsError:
		chatKey = t.Keys.SmsError
	case t.Keys.EschoolAppointmentNotification:
		chatKey = t.Keys.EschoolAppointmentNotification
	case t.Keys.TechnicalSupport:
		chatKey = t.Keys.TechnicalSupport
	default:
		return nil, fmt.Errorf("invalid chat key: %s", key)
	}

	return t.Contact(chatKey, phone, firstName, lastName)
}

// send makes the actual API request to Telegram
func (t *Telegram) Send(method string, parameters url.Values) (map[string]interface{}, error) {
	apiURL := fmt.Sprintf("%s%s/%s?%s", t.url, t.token, method, parameters.Encode())

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return result, nil
}

// setToken sets the Telegram bot token
func (t *Telegram) setToken() error {
	t.token = config.Get().Telegram.Token
	if t.token == "" {
		return fmt.Errorf("telegram API token not found")
	}
	return nil
}
