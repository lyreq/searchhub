package responses

type ListContentEnvelopeResponse struct {
	Data       []ListContentResponse `json:"data"`
	Pagination Pagination            `json:"pagination"`
	Sort       string                `json:"sort"`
	Filters    struct {
		Query string `json:"query,omitempty"`
		Type  string `json:"type,omitempty"`
	} `json:"filters"`
}

type ListContentResponse struct {
	ID          uint     `json:"id"`
	Title       string   `json:"title"`
	Type        string   `json:"type"`
	Score       float64  `json:"score"`
	PublishedAt string   `json:"published_at"`
	Tags        []string `json:"tags,omitempty"`
	Provider    string   `json:"provider"`
}

type ListCacheContentEntryResponse struct {
	Items []ListContentResponse `json:"items"`
	Total int64                 `json:"total"`
}

type ShowContentResponse struct {
	ID          uint     `json:"id"`
	Title       string   `json:"title"`
	Type        string   `json:"type"`
	Provider    string   `json:"provider"`
	PublishedAt string   `json:"published_at"`
	Tags        []string `json:"tags,omitempty"`

	Metrics struct {
		Views              *int `json:"views,omitempty"`
		Likes              *int `json:"likes,omitempty"`
		DurationSeconds    *int `json:"duration_seconds,omitempty"`
		ReadingTimeMinutes *int `json:"reading_time_minutes,omitempty"`
		Reactions          *int `json:"reactions,omitempty"`
		Comments           *int `json:"comments,omitempty"`
	} `json:"metrics"`

	Scores struct {
		Base            float64 `json:"base"`
		TypeCoefficient float64 `json:"type_coefficient"`
		Freshness       float64 `json:"freshness"`
		Engagement      float64 `json:"engagement"`
		Final           float64 `json:"final"`
	} `json:"scores"`

	RawPayload any `json:"raw_payload,omitempty"`
}
