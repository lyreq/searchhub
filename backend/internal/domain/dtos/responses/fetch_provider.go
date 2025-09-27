package responses

type ProviderRunResult struct {
	Provider string   `json:"provider"`
	Fetched  int      `json:"fetched"`
	Inserted int      `json:"inserted"`
	Updated  int      `json:"updated"`
	Skipped  int      `json:"skipped"`
	Errors   []string `json:"errors"`
}

type FetchResult struct {
	RunID     string              `json:"run_id"`
	Total     int                 `json:"total"`
	Inserted  int                 `json:"inserted"`
	Updated   int                 `json:"updated"`
	Skipped   int                 `json:"skipped"`
	Providers []ProviderRunResult `json:"providers"`
}
