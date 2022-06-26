package internal

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	gh "github.com/cli/go-gh"
	ghRepo "github.com/cli/go-gh/pkg/repository"
	"gopkg.in/h2non/gock.v1"
)

const MB_IN_BYTES = 1024 * 1024
const GB_IN_BYTES = 1024 * 1024 * 1024

var SORT_INPUT_TO_QUERY_MAP = map[string]string{
	"created-at": "created_at",
	"last-used":  "last_accessed_at",
	"size":       "size_in_bytes",
}

func GenerateQueryParams(branch string, limit int, key string, order string, sort string, page int) url.Values {
	query := url.Values{}
	if branch != "" {
		if strings.HasPrefix(branch, "refs/") {
			query.Add("ref", branch)
		} else {
			query.Add("ref", fmt.Sprintf("refs/heads/%s", branch))
		}
	}
	if limit != 30 {
		query.Add("per_page", strconv.Itoa(limit))
	}
	if key != "" {
		query.Add("key", key)
	}
	if order != "" {
		query.Add("direction", order)
	}
	if sort != "" {
		query.Add("sort", SORT_INPUT_TO_QUERY_MAP[sort])
	}
	if page > 1 {
		query.Add("page", strconv.Itoa(page))
	}

	return query
}

func GetRepo(r string) (ghRepo.Repository, error) {
	if r != "" {
		return ghRepo.Parse(r)
	}

	return gh.CurrentRepository()
}

func FormatCacheSize(size_in_bytes float64) string {
	if size_in_bytes < 1024 {
		return fmt.Sprintf("%.2f B", size_in_bytes)
	}

	if size_in_bytes < MB_IN_BYTES {
		return fmt.Sprintf("%.2f KB", size_in_bytes/1024)
	}

	if size_in_bytes < GB_IN_BYTES {
		return fmt.Sprintf("%.2f MB", size_in_bytes/MB_IN_BYTES)
	}

	return fmt.Sprintf("%.2f GB", size_in_bytes/GB_IN_BYTES)
}

func PrintPendingMocks(mocks []gock.Mock) string {
	paths := []string{}
	for _, mock := range mocks {
		paths = append(paths, mock.Request().URLStruct.String())
	}
	return fmt.Sprintf("%d unmatched mocks: %s", len(paths), strings.Join(paths, ", "))
}
