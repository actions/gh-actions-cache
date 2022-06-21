package cmd

import (
	"fmt"
	"net/url"
	"strconv"

	gh "github.com/cli/go-gh"
	ghRepo "github.com/cli/go-gh/pkg/repository"
)

const MB_IN_BYTES = 1024 * 1024
const GB_IN_BYTES = 1024 * 1024 * 1024

func generateQueryParams(branch string, limit int, key string, order string, sort string) url.Values {
	query := url.Values{}
	if branch != "" {
		query.Add("ref", branch)
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
		query.Add("sort", sort)
	}

	return query
}

func getRepo(r string) (ghRepo.Repository, error) {
	if r != "" {
		return ghRepo.Parse(r)
	}

	return gh.CurrentRepository()
}

func formatCacheSize(size_in_bytes float64) string {
	if size_in_bytes < 1024 {
		return fmt.Sprintf("%.2f Bytes", size_in_bytes)
	}

	if size_in_bytes < 1024*1024 {
		return fmt.Sprintf("%.2f KB", size_in_bytes/1024)
	}

	if size_in_bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", size_in_bytes/MB_IN_BYTES)
	}

	return fmt.Sprintf("%.2f GB", size_in_bytes/GB_IN_BYTES)
}
