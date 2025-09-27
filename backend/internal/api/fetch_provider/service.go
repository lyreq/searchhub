package fetch_provider

import (
	"context"
	"encoding/json"
	"lytemp/internal/domain/dtos/responses"
	"lytemp/internal/domain/models"
	"lytemp/internal/utils"
	"lytemp/pkg/provider"
	"time"
)

type Service struct {
	repository IRepository
}

func NewService(repository IRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) FetchAndStore(ctx context.Context, query string, pages, perPage int) (responses.FetchResult, error) {
	now := time.Now().UTC()
	runID := now.Format("20060102-150405")

	adapters, err := provider.BuildProviders()
	if err != nil {
		return responses.FetchResult{}, err
	}

	res := responses.FetchResult{
		RunID: runID,
	}

	for _, p := range adapters {
		pres := responses.ProviderRunResult{Provider: p.Name()}

		for page := 1; page <= pages; page++ {
			pageData, err := p.Fetch(ctx, page, perPage, query)
			if err != nil {
				pres.Errors = append(pres.Errors, err.Error())
				break
			}
			pres.Fetched += len(pageData.Items)

			for _, it := range pageData.Items {
				// skorları hesapla
				b := utils.BaseScore(it)
				tc := utils.TypeCoef(it.NormalizedType)
				f := utils.FreshnessScore(it.PublishedAt, now)
				e := utils.EngagementScore(it)

				final := (b * tc) + f + e

				content := models.Content{
					Provider:        it.Provider,
					ProviderItemID:  it.ProviderItemID,
					Title:           it.Title,
					NormalizedType:  models.ContentType(it.NormalizedType),
					PublishedAt:     it.PublishedAt,
					BaseScore:       utils.Round2(b),
					FreshnessScore:  utils.Round2(f),
					EngagementScore: utils.Round2(e),
					FinalScore:      utils.Round2(final),
					IndexedAt:       now,
				}

				if it.Views != nil {
					content.Views = it.Views
				}
				if it.Likes != nil {
					content.Likes = it.Likes
				}
				if it.DurationSeconds != nil {
					content.DurationSeconds = it.DurationSeconds
				}
				if it.ReadingTimeMinutes != nil {
					content.ReadingTimeMinutes = it.ReadingTimeMinutes
				}
				if it.Reactions != nil {
					content.Reactions = it.Reactions
				}
				if it.Comments != nil {
					content.Comments = it.Comments
				}
				if len(it.Tags) > 0 {
					if bts, err := json.Marshal(it.Tags); err == nil {
						content.Tags = bts
					}
				}
				if it.Raw != nil {
					if bts, err := json.Marshal(it.Raw); err == nil {
						content.RawPayload = bts
					}
				}

				created, err := s.repository.UpsertContent(ctx, &content)
				if err != nil {
					pres.Errors = append(pres.Errors, err.Error())
					continue
				}
				if created {
					pres.Inserted++
				} else {
					pres.Updated++
				}
			}
		}

		status := models.ProviderStatusOK
		var lastErrStr *string
		if len(pres.Errors) > 0 {
			status = models.ProviderStatusDegraded
			msg := pres.Errors[len(pres.Errors)-1]
			lastErrStr = &msg
		}
		_ = s.repository.UpsertProviderSync(ctx, p.Name(), utils.DetectFormat(p.Name()), 60, status, now, lastErrStr)

		res.Providers = append(res.Providers, pres)
		res.Total += pres.Fetched
		res.Inserted += pres.Inserted
		res.Updated += pres.Updated
		res.Skipped += pres.Skipped
	}

	return res, nil
}
