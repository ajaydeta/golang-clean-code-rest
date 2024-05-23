package domain

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type Filter struct {
	Limit  int64 `cache_key:"limit"`
	Offset int64 `cache_key:"offset"`
	Search string
	Sort   string `cache_key:"sort"`
	Order  string `cache_key:"order"`
}

func (f Filter) HasSearch() bool {
	return len(strings.TrimSpace(f.Search)) > 0
}

func (f Filter) GetSort() string {
	if "DESC" == strings.ToUpper(f.Sort) {
		return "DESC"
	}

	return "ASC"
}

func (f Filter) GetSortAndPaginationWithDefaultQuery(db *gorm.DB, defaultOrder string, columnMapFilter map[string]string) *gorm.DB {

	orderQuery := defaultOrder
	if f.Order != "" {
		if column, found := columnMapFilter[strings.ToLower(f.Order)]; found {
			orderQuery = fmt.Sprintf("%s %s", column, f.GetSort())
		}
	}

	if f.Limit == 0 {
		f.Limit = 10
	}

	return db.Order(orderQuery).Limit(int(f.Limit)).Offset(int(f.Offset))
}
