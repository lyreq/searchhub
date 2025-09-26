package notification

import (
	"lytemp/config"
	"lytemp/pkg/onesignal"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	onesignal     *onesignal.OneSignalService
}
type NotificationJob struct {
	PlayerID string                 `json:"player_id"`
	Title    string                 `json:"title"`
	Content  string                 `json:"content"`
	Data     map[string]interface{} `json:"data"`
}

func New(onesignal *onesignal.OneSignalService) *Service {
	return &Service{onesignal: onesignal}
}

func (s *Service) EnqueueNotificationJob(ctx context.Context, job NotificationJob) error {
	cfg := config.Get()
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Pass,
		DB:       0,
	})

	payload, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return rdb.RPush(ctx, "onesignal:queue", payload).Err()
}

func (s *Service) SendCustomNotification(playerID, title, message string, data map[string]interface{}) error {
	//playerID = s.ContactFilter.FilterNotificationContacts(playerID)
	//if playerID == "" {
	//	return nil
	//}

	log.Println("playerID", playerID)

	notification := s.onesignal.
		SetTitle(title).
		SetMessage(message).
		SetPlayer(playerID).
		MustHavePlayers()

	// Eğer özel data varsa ekle
	for key, value := range data {
		notification.SetAdditional(key, value)
	}

	ok, err := notification.Send()
	if err != nil || !ok {
		log.Println("onesignal bildirim gönderilemedi: %w", err)
		return fmt.Errorf("onesignal bildirim gönderilemedi: %w", err)
	}
	return nil
}
