package provider

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ----- wire types matching the XML mock -----
//
// <feed>
//   <items>
//     <item>...</item>
//   </items>
//   <meta>...</meta>
// </feed>

type xmlFeed struct {
	XMLName xml.Name `xml:"feed"`
	Items   xmlItems `xml:"items"`
	Meta    xmlMeta  `xml:"meta"`
}

type xmlItems struct {
	List []xmlItem `xml:"item"`
}

type xmlItem struct {
	ID              string   `xml:"id"`
	Headline        string   `xml:"headline"`
	Type            string   `xml:"type"` // "video" | "article"
	Stats           xmlStats `xml:"stats"`
	PublicationDate string   `xml:"publication_date"` // "YYYY-MM-DD"
	Categories      xmlCats  `xml:"categories"`
}

type xmlStats struct {
	// video fields
	Views    *int   `xml:"views"`
	Likes    *int   `xml:"likes"`
	Duration string `xml:"duration"` // "MM:SS" or "HH:MM:SS"

	// article fields
	ReadingTime *int `xml:"reading_time"` // minutes
	Reactions   *int `xml:"reactions"`
	Comments    *int `xml:"comments"`
}

type xmlCats struct {
	Category []string `xml:"category"`
}

type xmlMeta struct {
	TotalCount   int `xml:"total_count"`
	CurrentPage  int `xml:"current_page"`
	ItemsPerPage int `xml:"items_per_page"`
}

// ----- adapter -----

type XMLProvider struct {
	name   string
	url    string
	client *http.Client
}

func NewXMLProvider(url string, client *http.Client) *XMLProvider {
	if client == nil {
		client = HTTPClient()
	}
	// default: 60 RPM, burst 5
	SetRateLimit("xml_provider", 60, 5)

	return &XMLProvider{
		name:   "xml_provider",
		url:    url,
		client: client,
	}
}
func (p *XMLProvider) Name() string { return p.name }

func (p *XMLProvider) Fetch(ctx context.Context, page, perPage int, query string) (Page, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.url, nil)
	if err != nil {
		return Page{}, err
	}
	req.Header.Set("User-Agent", "LyTemp-ContentSearch/1.0 (+xml)")
	resp, err := Do(ctx, p.name, req, p.client)
	if err != nil {
		return Page{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return Page{}, fmt.Errorf("xml provider: unexpected status %d: %s", resp.StatusCode, string(b))
	}

	var feed xmlFeed
	dec := xml.NewDecoder(resp.Body)
	dec.Strict = false
	if err := dec.Decode(&feed); err != nil {
		return Page{}, err
	}

	out := Page{
		Items:   make([]Item, 0, len(feed.Items.List)),
		Total:   feed.Meta.TotalCount,
		Page:    feed.Meta.CurrentPage,
		PerPage: feed.Meta.ItemsPerPage,
	}

	for _, it := range feed.Items.List {
		// publication_date: "YYYY-MM-DD" -> UTC midnight
		pub := parseYYYYMMDD(it.PublicationDate)

		normalizedType := strings.ToLower(it.Type)

		var viewsPtr, likesPtr, durPtr *int
		var rtPtr, rePtr, cmPtr *int

		// video stats (if present)
		if it.Stats.Views != nil {
			viewsPtr = it.Stats.Views
		}
		if it.Stats.Likes != nil {
			likesPtr = it.Stats.Likes
		}
		if it.Stats.Duration != "" {
			d := parseDurationToSeconds(it.Stats.Duration)
			durPtr = &d
		}

		// article stats (if present)
		if it.Stats.ReadingTime != nil {
			rtPtr = it.Stats.ReadingTime
		}
		if it.Stats.Reactions != nil {
			rePtr = it.Stats.Reactions
		}
		if it.Stats.Comments != nil {
			cmPtr = it.Stats.Comments
		}

		item := Item{
			Provider:           p.name,
			ProviderItemID:     it.ID,
			Title:              it.Headline,
			NormalizedType:     normalizedType, // "video" | "article"
			Views:              viewsPtr,
			Likes:              likesPtr,
			DurationSeconds:    durPtr,
			ReadingTimeMinutes: rtPtr,
			Reactions:          rePtr,
			Comments:           cmPtr,
			PublishedAt:        pub,
			Tags:               it.Categories.Category,
			Raw:                it,
		}
		out.Items = append(out.Items, item)
	}

	// Respect caller pagination overrides if provided.
	if page > 0 {
		out.Page = page
	}
	if perPage > 0 {
		out.PerPage = perPage
	}
	return out, nil
}

func parseYYYYMMDD(s string) time.Time {
	if s == "" {
		return time.Now().UTC()
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		// try more formats, then fallback
		if n, e := strconv.ParseInt(s, 10, 64); e == nil && n > 0 {
			return time.Unix(n, 0).UTC()
		}
		return time.Now().UTC()
	}
	// midnight UTC
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
