package meilisearch

import (
	"lytemp/internal/utils/algorithm"
	"time"

	"github.com/google/uuid"
)

// func (c *client) Log(index Index, doc map[string]interface{}) {
// 	if err := c.Add(string(index), doc); err != nil {
// 		log.Println(err)
// 	}
// }

func (c *client) LogUserActivity(message string, userId uint) {
	doc := map[string]interface{}{
		"id":             uuid.New().String(),
		"message":        message,
		"process_name":   algorithm.GetFunctionName(3),
		"user_id":        userId,
		"application_id": 1,
		"created_at":     time.Now(),
		"created_at_ts":  time.Now().Unix(), // meili search sadece unix timestamp ile filtreleme yapabiliyor
	}

	c.Log(UserActivityLog, doc)
}

func (c *client) LogUserErrLog(message string, userId uint) {
	doc := map[string]interface{}{
		"id":             uuid.New().String(),
		"message":        message,
		"process_name":   algorithm.GetFunctionName(3),
		"user_id":        userId,
		"application_id": 1,
		"created_at":     time.Now(),
		"created_at_ts":  time.Now().Unix(), // meili search sadece unix timestamp ile filtreleme yapabiliyor
	}

	c.Log(UserErrLog, doc)
}

func (c *client) LogSystemErrLog(message string) {
	doc := map[string]interface{}{
		"id":             uuid.New().String(),
		"message":        message,
		"process_name":   algorithm.GetFunctionName(3),
		"application_id": 1,
		"created_at":     time.Now(),
		"created_at_ts":  time.Now().Unix(), // meili search sadece unix timestamp ile filtreleme yapabiliyor
	}

	c.Log(SystemErrLog, doc)
}
