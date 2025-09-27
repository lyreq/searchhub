package utils

import (
	"math"
	"testing"
	"time"

	"lytemp/pkg/provider"

	"github.com/stretchr/testify/require"
)

func TestBaseScore_Video(t *testing.T) {
	views := 15000
	likes := 1200
	it := provider.Item{NormalizedType: "video", Views: &views, Likes: &likes}
	got := BaseScore(it)
	require.InDelta(t, 27.0, got, 0.0001) // 15 + 12
}

func TestBaseScore_Article(t *testing.T) {
	rt := 8
	re := 450
	it := provider.Item{NormalizedType: "article", ReadingTimeMinutes: &rt, Reactions: &re}
	got := BaseScore(it)
	require.InDelta(t, 17.0, got, 0.0001) // 8 + 9
}

func TestTypeCoef(t *testing.T) {
	require.Equal(t, 1.5, TypeCoef("video"))
	require.Equal(t, 1.0, TypeCoef("article"))
	require.Equal(t, 1.0, TypeCoef("unknown"))
}

func TestFreshnessScore(t *testing.T) {
	now := time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC)
	require.Equal(t, 5.0, FreshnessScore(now.Add(-6*24*time.Hour), now))
	require.Equal(t, 3.0, FreshnessScore(now.Add(-20*24*time.Hour), now))
	require.Equal(t, 1.0, FreshnessScore(now.Add(-70*24*time.Hour), now))
	require.Equal(t, 0.0, FreshnessScore(now.Add(-200*24*time.Hour), now))
}

func TestEngagementScore_Video(t *testing.T) {
	views := 25000
	likes := 2100
	it := provider.Item{NormalizedType: "video", Views: &views, Likes: &likes}
	got := EngagementScore(it)
	require.InDelta(t, 0.84, got, 1e-9) // (2100/25000)*10
}

func TestEngagementScore_Article(t *testing.T) {
	rt := 8
	re := 450
	it := provider.Item{NormalizedType: "article", ReadingTimeMinutes: &rt, Reactions: &re}
	got := EngagementScore(it)
	require.True(t, math.Abs(got-281.25) < 1e-9) // (450/8)*5
}
