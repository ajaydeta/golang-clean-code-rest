package shared

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	AccessTokenSubject  = "access_token"
	RefreshTokenSubject = "refresh_token"

	PaymentTypeTransferBank = "1_transfer_bank"
	PaymentTypeSupermarket  = "2_supermarket"

	AccessTokenDuration  = time.Hour * 24
	RefreshTokenDuration = time.Hour * 24 * 7
	CacheTtl             = time.Minute * 30
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

func GetCacheKey(prefix string, i interface{}) string {
	toString := fmt.Sprintf("%+v", i)
	noWhiteSpace := strings.ReplaceAll(toString, " ", "_")
	noOpenBracket := strings.ReplaceAll(noWhiteSpace, "{", "")
	noCloseBracket := strings.ReplaceAll(noOpenBracket, "}", "")
	result := prefix + noCloseBracket

	return result
}
