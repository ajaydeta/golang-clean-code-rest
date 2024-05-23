package domain

import (
	"fmt"
	"net/url"
	"strings"
)

type Filter struct {
	Limit    int64 `cache_key:"limit"`
	Offset   int64 `cache_key:"offset"`
	Search   string
	Sort     string `cache_key:"sort"`
	Order    string `cache_key:"order"`
	UrlQuery url.Values
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

func (f Filter) GetOrder(orderBy map[string]string, args []interface{}) (string, []interface{}) {
	var orderByText = ""
	if _, ok := orderBy[f.Order]; ok {
		orderByText = " ORDER BY " + orderBy[f.Order]
		orderByText += " " + f.GetSort()
	}
	if f.Limit < 1 {
		f.Limit = 10
	}
	orderByText += " LIMIT ? OFFSET ?"
	args = append(args, f.Limit, f.Offset)
	return orderByText, args
}

func (f Filter) GetSortAndPaginationWithDefaultQuery(query, defaultOrder string, columnMapFilter map[string]string, args []interface{}) (string, []interface{}) {
	var (
		orderQuery = "ORDER BY " + defaultOrder
	)

	if f.Order != "" {
		if column, found := columnMapFilter[strings.ToLower(f.Order)]; found {
			orderQuery = fmt.Sprintf("ORDER BY %s %s", column, f.GetSort())
		}
	}

	limitQuery := "LIMIT ? OFFSET ?"
	if f.Limit == 0 {
		f.Limit = 10
	}
	args = append(args, f.Limit, f.Offset)

	query = strings.Join([]string{query, orderQuery, limitQuery}, " ")

	return query, args
}
