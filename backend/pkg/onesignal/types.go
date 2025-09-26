package onesignal

// NotificationRequest OneSignal API'sine gönderilecek bildirim isteği
type NotificationRequest struct {
	AppID            string                 `json:"app_id"`
	Contents         map[string]string      `json:"contents,omitempty"`
	Headings         map[string]string      `json:"headings,omitempty"`
	IncludePlayerIDs []string               `json:"include_player_ids,omitempty"`
	IncludedSegments []string               `json:"included_segments,omitempty"`
	Data             map[string]interface{} `json:"data,omitempty"`
}

// NotificationResponse OneSignal API'sinden dönen yanıt
type NotificationResponse struct {
	ID      string   `json:"id"`
	Errors  []string `json:"errors,omitempty"`
	Success bool     `json:"success"`
}

// OneSignalService OneSignal bildirim servisi
type OneSignalService struct {
	appID           string
	apiKey          string
	data            map[string]interface{}
	additionals     map[string]interface{}
	title           string
	message         string
	players         []string
	mustHavePlayers bool
}
