package onesignal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// NewOneSignalService yeni bir OneSignal servisi oluşturur
func NewOneSignalService(appID, apiKey string) *OneSignalService {
	return &OneSignalService{
		appID:           appID,
		apiKey:          apiKey,
		data:            make(map[string]interface{}),
		additionals:     make(map[string]interface{}),
		title:           "Giriş Doğrulama",
		message:         "",
		players:         make([]string, 0),
		mustHavePlayers: false,
	}
}

// prepare bildirim verilerini hazırlar
func (o *OneSignalService) prepare() {
	o.data["app_id"] = o.appID
	o.data["contents"] = map[string]string{"en": o.message}
	o.data["headings"] = map[string]string{"en": o.title}

	if len(o.players) > 0 {
		o.data["include_player_ids"] = o.players
	} else {
		o.data["included_segments"] = []string{"All"}
	}

	if len(o.additionals) > 0 {
		o.data["data"] = o.additionals
	}
}

// sendToOneSignal OneSignal API'sine bildirim gönderir
func (o *OneSignalService) sendToOneSignal() (bool, error) {
	o.prepare()

	jsonData, err := json.Marshal(o.data)
	if err != nil {
		log.Println("JSON marshal hatası: %v", err)
		return false, fmt.Errorf("JSON marshal hatası: %v", err)
	}

	req, err := http.NewRequest("POST", "https://onesignal.com/api/v1/notifications", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("HTTP request oluşturma hatası: %v", err)
		return false, fmt.Errorf("HTTP request oluşturma hatası: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+o.apiKey)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("HTTP istek hatası: %v", err)
		return false, fmt.Errorf("HTTP istek hatası: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("yanıt okuma hatası: %v", err)
		return false, fmt.Errorf("yanıt okuma hatası: %v", err)
	}

	var response NotificationResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Println("JSON unmarshal hatası: %v", err)
		return false, fmt.Errorf("JSON unmarshal hatası: %v", err)
	}

	return len(response.Errors) == 0, nil
}

// Send bildirimi gönderir
func (o *OneSignalService) Send() (bool, error) {
	if !o.mustHavePlayers {
		return o.sendToOneSignal()
	} else if o.mustHavePlayers && len(o.players) > 0 {
		return o.sendToOneSignal()
	}

	return true, nil
}

// SetAdditional ek veri ekler
func (o *OneSignalService) SetAdditional(key string, value interface{}) *OneSignalService {
	o.additionals[key] = value
	return o
}

// SetPlayer tek bir player ID ekler
func (o *OneSignalService) SetPlayer(player string) *OneSignalService {
	o.players = append(o.players, player)
	return o
}

// SetPlayers birden fazla player ID ekler
func (o *OneSignalService) SetPlayers(players []string) *OneSignalService {
	o.players = append(o.players, players...)
	return o
}

// SetTitle bildirim başlığını ayarlar
func (o *OneSignalService) SetTitle(title string) *OneSignalService {
	o.title = title
	return o
}

// SetMessage bildirim mesajını ayarlar
func (o *OneSignalService) SetMessage(message string) *OneSignalService {
	o.message = message
	return o
}

// MustHavePlayers player ID'lerin zorunlu olmasını sağlar
func (o *OneSignalService) MustHavePlayers() *OneSignalService {
	o.mustHavePlayers = true
	return o
}

// Reset servisi sıfırlar
func (o *OneSignalService) Reset() *OneSignalService {
	o.data = make(map[string]interface{})
	o.additionals = make(map[string]interface{})
	o.title = "Giriş Doğrulama"
	o.message = ""
	o.players = make([]string, 0)
	o.mustHavePlayers = false
	return o
}
