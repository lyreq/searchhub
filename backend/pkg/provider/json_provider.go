package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ----- wire types matching the JSON mock -----

type jsonPagination struct {
	Total   int `json:"total"`
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

type jsonMetrics struct {
	Views    int    `json:"views"`
	Likes    int    `json:"likes"`
	Duration string `json:"duration"` // "MM:SS" or "HH:MM:SS"
}

type jsonContent struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Type        string      `json:"type"` // "video"
	Metrics     jsonMetrics `json:"metrics"`
	PublishedAt string      `json:"published_at"` // ISO8601
	Tags        []string    `json:"tags"`
}

type jsonEnvelope struct {
	Contents   []jsonContent  `json:"contents"`
	Pagination jsonPagination `json:"pagination"`
}

// ----- adapter -----

type JSONProvider struct {
	name   string
	url    string
	client *http.Client
}

func NewJSONProvider(url string, client *http.Client) *JSONProvider {
	if client == nil {
		client = HTTPClient()
	}
	// default: 60 RPM, burst 5
	SetRateLimit("json_provider", 60, 5)

	return &JSONProvider{
		name:   "json_provider",
		url:    url,
		client: client,
	}
}

func (p *JSONProvider) Name() string { return p.name }

func (p *JSONProvider) Fetch(ctx context.Context, page, perPage int, query string) (Page, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.url, nil)
	if err != nil {
		return Page{}, err
	}
	req.Header.Set("User-Agent", "LyTemp-ContentSearch/1.0 (+json)")
	resp, err := Do(ctx, p.name, req, p.client)
	if err != nil {
		return Page{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return Page{}, fmt.Errorf("json provider: unexpected status %d: %s", resp.StatusCode, string(b))
	}

	var env jsonEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return Page{}, err
	}

	out := Page{
		Items:   make([]Item, 0, len(env.Contents)),
		Total:   env.Pagination.Total,
		Page:    env.Pagination.Page,
		PerPage: env.Pagination.PerPage,
	}

	// Normalize each content.
	for _, c := range env.Contents {
		published, err := time.Parse(time.RFC3339, c.PublishedAt)
		if err != nil {
			// be forgiving: try common fallback
			if t2, e2 := time.Parse("2006-01-02T15:04:05Z07:00", c.PublishedAt); e2 == nil {
				published = t2
			} else {
				published = time.Now().UTC()
			}
		}

		// Duration "MM:SS" or "HH:MM:SS" -> seconds
		durSec := parseDurationToSeconds(c.Metrics.Duration)

		views := c.Metrics.Views
		likes := c.Metrics.Likes

		item := Item{
			Provider:        p.name,
			ProviderItemID:  c.ID,
			Title:           c.Title,
			NormalizedType:  strings.ToLower(c.Type), // expect "video"
			Views:           ptrInt(views),
			Likes:           ptrInt(likes),
			DurationSeconds: ptrInt(durSec),
			PublishedAt:     published.UTC(),
			Tags:            c.Tags,
			Raw:             c, // keep raw object for debug
		}
		out.Items = append(out.Items, item)
	}

	// If caller passed explicit page/perPage, prefer them in the wrapper (data stays same).
	if page > 0 {
		out.Page = page
	}
	if perPage > 0 {
		out.PerPage = perPage
	}
	return out, nil
}

// parseDurationToSeconds supports "MM:SS" and "HH:MM:SS".
func parseDurationToSeconds(s string) int {
	if s == "" {
		return 0
	}
	parts := strings.Split(s, ":")
	toInt := func(x string) int {
		n, _ := strconv.Atoi(x)
		return n
	}
	switch len(parts) {
	case 2: // MM:SS
		m := toInt(parts[0])
		sec := toInt(parts[1])
		return m*60 + sec
	case 3: // HH:MM:SS
		h := toInt(parts[0])
		m := toInt(parts[1])
		sec := toInt(parts[2])
		return h*3600 + m*60 + sec
	default:
		// try as seconds number
		n, _ := strconv.Atoi(s)
		return n
	}
}

func ptrInt(v int) *int { return &v }
