package requests


type ListContentRequest struct {
	Query   string `query:"query" validate:"omitempty"`
	Type    string `query:"type" validate:"omitempty,oneof=video article"`
	Sort    string `query:"sort" validate:"omitempty,oneof=score relevance recent popularity"`
	Page    int    `query:"page" validate:"omitempty,min=1"`
	PerPage int    `query:"per_page" validate:"omitempty,min=1,max=100"`
}



type ContentShowRequest struct {
	ID         uint  `param:"id" validate:"required,gt=0"`
	IncludeRaw *bool `query:"include_raw" validate:"omitempty"`
}