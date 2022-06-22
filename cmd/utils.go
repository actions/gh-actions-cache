package cmd

import (
	"fmt"
	"math"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	gh "github.com/cli/go-gh"
	ghRepo "github.com/cli/go-gh/pkg/repository"
	"github.com/moby/term"
)

const MB_IN_BYTES = 1024 * 1024
const GB_IN_BYTES = 1024 * 1024 * 1024

var SORT_INPUT_TO_QUERY_MAP = map[string]string{
	"created-at": "created_at",
	"last-used":  "last_accessed_at",
	"size":       "size_in_bytes",
}

func generateQueryParams(branch string, limit int, key string, order string, sort string) url.Values {
	query := url.Values{}
	if branch != "" {
		if strings.HasPrefix(branch, "refs/"){
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

func prettyPrintCacheList(caches []cacheInfo) {
	numberOfCaches := len(caches)
	for _, cache := range caches {
		var cacheKey string = trimCacheKeyBasedOnWindowSize(cache.Key)

		var formattedRow string = getFormattedCacheInfo(cacheKey, cache)
		fmt.Println(formattedRow)
	}
	if numberOfCaches > 30 {
		fmt.Println("...and " + strconv.Itoa(numberOfCaches-30) + " more\n\n")
	}
	fmt.Println()
}

func trimCacheKeyBasedOnWindowSize(key string) string {
	var cacheKey string
	fd := os.Stdin.Fd()
	ws, _ := term.GetWinsize(fd)
	var cacheKeyWidth int = int(math.Max(10, float64(ws.Width)-60))
	if len(key) > int(cacheKeyWidth) {
		cacheKey = key[:cacheKeyWidth] + "..."
	} else {
		cacheKey = key
		for i := 0; i < 5; i++ {
			key += " "
		}
	}
	return cacheKey
}

func lastAccessedHour(lastAccessedAt string) string {
	var now time.Time = time.Now()
	lastAccessedTime, _ := time.Parse(time.RFC3339, lastAccessedAt)

	diff := now.Sub(lastAccessedTime)
	lastAccessedHourStr := "Used " + strconv.FormatFloat(diff.Hours(), 'f', 0, 64)
	if diff.Hours() < 2 {
		lastAccessedHourStr += " hour ago"
	} else {
		lastAccessedHourStr += " hours ago"
	}
	return lastAccessedHourStr
}

func getFormattedCacheInfo(cacheKey string, cache cacheInfo) string {
	return "  " + cacheKey + "     [" + formatCacheSize(cache.Size) + "]     " + cache.Ref[11:] + "     " + lastAccessedHour(cache.LastAccessedAt)
}
