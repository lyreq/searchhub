package requests

type FetchProviderRequest struct {
	Query   string `query:"query" validate:"omitempty"`
	Pages   int    `query:"pages" validate:"omitempty,min=1,max=20"`
	PerPage int    `query:"per_page" validate:"omitempty,min=1,max=100"`
}
