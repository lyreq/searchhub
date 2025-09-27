package contents

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"lytemp/internal/domain/dtos/responses"
	"lytemp/internal/utils"
)

type Service struct {
	r IRepository

	cacheEnabled bool
	cacheGet     func(ctx context.Context, key string) (string, error)
	cacheSet     func(ctx context.Context, key string, value string, ttl time.Duration) error
	cacheTTL     time.Duration
}

func NewService(r IRepository) *Service { return &Service{r: r} }

func (s *Service) WithCache(
	get func(ctx context.Context, key string) (string, error),
	set func(ctx context.Context, key string, value string, ttl time.Duration) error,
	ttl time.Duration,
) *Service {
	if ttl <= 0 {
		ttl = 60 * time.Second
	}
	s.cacheEnabled = true
	s.cacheGet = get
	s.cacheSet = set
	s.cacheTTL = ttl
	return s
}
func (s *Service) List(ctx context.Context, Query string, Type string, Sort string, Page int, PerPage int) (responses.ListContentEnvelopeResponse, error) {
	var out responses.ListContentEnvelopeResponse
	if Page < 1 {
		Page = 1
	}
	if PerPage <= 0 {
		PerPage = 20
	}
	if PerPage > 100 {
		PerPage = 100
	}

	cacheKey := s.makeCacheKey(Query, Type, Sort, Page, PerPage)

	if s.cacheEnabled && s.cacheGet != nil {
		if raw, err := s.cacheGet(ctx, cacheKey); err == nil && raw != "" {
			var cached struct {
				Items []responses.ListContentResponse `json:"items"`
				Total int64                           `json:"total"`
			}
			if json.Unmarshal([]byte(raw), &cached) == nil && len(cached.Items) > 0 {
				out.Data = cached.Items
				out.Pagination = responses.Pagination{
					Total:   cached.Total,
					Page:    utils.DefInt(Page, 1),
					PerPage: utils.DefInt(PerPage, 20),
				}
				out.Sort = Sort
				out.Filters.Query = Query
				out.Filters.Type = Type
				return out, nil
			}
		}
	}

	items, total, err := s.r.Search(ctx, Query, Type, Sort, Page, PerPage)
	if err != nil {
		return out, err
	}

	list := make([]responses.ListContentResponse, 0, len(items))
	for _, c := range items {
		var tags []string
		if len(c.Tags) > 0 {
			_ = json.Unmarshal(c.Tags, &tags)
		}
		list = append(list, responses.ListContentResponse{
			ID:          c.ID,
			Title:       c.Title,
			Type:        string(c.NormalizedType),
			Score:       c.FinalScore,
			PublishedAt: c.PublishedAt.UTC().Format(time.RFC3339),
			Tags:        tags,
			Provider:    c.Provider,
		})
	}

	out.Data = list
	out.Pagination = responses.Pagination{
		Total:   total,
		Page:    utils.DefInt(Page, 1),
		PerPage: utils.DefInt(PerPage, 20),
	}
	out.Sort = Sort
	out.Filters.Query = Query
	out.Filters.Type = Type

	if s.cacheEnabled && s.cacheSet != nil && len(list) > 0 {
		payload, _ := json.Marshal(struct {
			Items []responses.ListContentResponse `json:"items"`
			Total int64                           `json:"total"`
		}{Items: list, Total: total})
		_ = s.cacheSet(ctx, cacheKey, string(payload), s.cacheTTL)
	}

	return out, nil
}

func (s *Service) Detail(ctx context.Context, id uint, includeRaw bool) (responses.ShowContentResponse, error) {
	var out responses.ShowContentResponse

	c, err := s.r.GetByID(ctx, id)
	if err != nil {
		return out, err
	}

	var tags []string
	if len(c.Tags) > 0 {
		_ = json.Unmarshal(c.Tags, &tags)
	}

	out.ID = c.ID
	out.Title = c.Title
	out.Type = string(c.NormalizedType)
	out.Provider = c.Provider
	out.PublishedAt = c.PublishedAt.UTC().Format(time.RFC3339)
	out.Tags = tags

	out.Metrics.Views = c.Views
	out.Metrics.Likes = c.Likes
	out.Metrics.DurationSeconds = c.DurationSeconds
	out.Metrics.ReadingTimeMinutes = c.ReadingTimeMinutes
	out.Metrics.Reactions = c.Reactions
	out.Metrics.Comments = c.Comments

	out.Scores.Base = c.BaseScore
	out.Scores.Freshness = c.FreshnessScore
	out.Scores.Engagement = c.EngagementScore
	out.Scores.Final = c.FinalScore
	out.Scores.TypeCoefficient = utils.TypeCoefFor(string(c.NormalizedType))

	if includeRaw && len(c.RawPayload) > 0 {
		var raw any
		if err := json.Unmarshal(c.RawPayload, &raw); err == nil {
			out.RawPayload = raw
		}
	}
	return out, nil
}

func (s *Service) makeCacheKey(Query string, Type string, Sort string, Page int, PerPage int) string {
	query := strings.TrimSpace(strings.ToLower(Query))
	typ := strings.ToLower(Type)
	sort := strings.ToLower(Sort)
	qEsc := url.QueryEscape(query)

	return strings.Join([]string{
		"contents:list:v1",
		"q=" + qEsc,
		"t=" + typ,
		"s=" + sort,
		"p=" + utils.StrconvI(Page, 1),
		"pp=" + utils.StrconvI(PerPage, 20),
	}, "|")
}
