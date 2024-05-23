package dto

type Pagination struct {
	Page      int64  `json:"page,omitempty"`
	PerPage   int64  `json:"per_page,omitempty"`
	TotalData int64  `json:"total_data,omitempty"`
	Sort      string `json:"sort,omitempty"`
	Order     string `json:"order,omitempty"`
}
