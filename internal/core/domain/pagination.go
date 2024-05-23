package domain

type Pagination struct {
	Page      int64
	PerPage   int64
	TotalData int64
	Sort      string
	Order     string
}
