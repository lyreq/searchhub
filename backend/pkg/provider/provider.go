package provider

import (
	"context"
	"time"
)

// Normalized item all providers return.
type Item struct {
	Provider          string
	ProviderItemID    string
	Title             string
	NormalizedType    string // "video" | "article"
	Views             *int
	Likes             *int
	DurationSeconds   *int
	ReadingTimeMinutes *int
	Reactions         *int
	Comments          *int
	PublishedAt       time.Time
	Tags              []string
	Raw               any // original payload (map or struct) for debugging
}

// Page is a normalized pagination wrapper.
type Page struct {
	Items   []Item
	Total   int
	Page    int
	PerPage int
}

// Provider is the common interface all adapters implement.
type Provider interface {
	Name() string
	Fetch(ctx context.Context, page, perPage int, query string) (Page, error)
}
