package provider_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"lytemp/pkg/provider"

	"github.com/stretchr/testify/require"
)

func TestJSONProvider_Fetch_Normalize(t *testing.T) {
	// Mock JSON envelope
	payload := map[string]any{
		"contents": []any{
			map[string]any{
				"id":    "v2",
				"title": "Advanced Go Concurrency Patterns",
				"type":  "video",
				"metrics": map[string]any{
					"views":    25000,
					"likes":    2100,
					"duration": "22:45",
				},
				"published_at": "2024-03-14T15:30:00Z",
				"tags":         []string{"programming", "advanced", "concurrency"},
			},
		},
		"pagination": map[string]any{
			"total":    4,
			"page":     1,
			"per_page": 10,
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(payload)
	}))
	defer srv.Close()

	p := provider.NewJSONProvider(srv.URL, provider.HTTPClient())
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	page, err := p.Fetch(ctx, 1, 10, "")
	require.NoError(t, err)
	require.Equal(t, 4, page.Total)
	require.Len(t, page.Items, 1)

	it := page.Items[0]
	require.Equal(t, "json_provider", it.Provider)
	require.Equal(t, "v2", it.ProviderItemID)
	require.Equal(t, "Advanced Go Concurrency Patterns", it.Title)
	require.Equal(t, "video", it.NormalizedType)
	require.NotNil(t, it.Views)
	require.NotNil(t, it.Likes)
	require.NotNil(t, it.DurationSeconds)
	require.Equal(t, 1365, *it.DurationSeconds) 
	require.Equal(t, []string{"programming", "advanced", "concurrency"}, it.Tags)
}
