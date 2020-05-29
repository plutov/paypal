package paypal

type ListParams struct{
	Page          string `json:"page,omitempty"`           //Default: 0.
	PageSize      string `json:"page_size,omitempty"`      //Default: 10.
	TotalRequired string `json:"total_required,omitempty"` //Default: no.
}

type ListResponse struct{
	TotalItems string        `json:"total_items,omitempty"`
	TotalPages string        `json:"total_pages,omitempty"`
	Links      []Link        `json:"links,omitempty"`
}
