package shared

import (
	"strconv"
	"strings"
	"time"
)

const (
	AccessTokenSubject  = "access_token"
	RefreshTokenSubject = "refresh_token"

	AccessTokenDuration  = time.Hour * 24
	RefreshTokenDuration = time.Hour * 24 * 7
)

func GetPageAndPerPage(limit, offset int64) (int64, int64) {
	page := offset
	perPage := limit
	if perPage != 0 {
		page = (offset + limit) / limit
	}
	return page, perPage
}

func StringToInt64(str string, defaultValue int64) int64 {
	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return defaultValue
	}

	return value
}

func SplitStringBySeparator(in string, sep string) []string {
	var sepStr []string
	for _, s := range strings.Split(in, sep) {
		s = strings.TrimSpace(s)
		if s != "" {
			sepStr = append(sepStr, s)
		}
	}

	return sepStr
}
