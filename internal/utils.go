package internal

import (
	"fmt"
	"math"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/TwiN/go-color"
	"github.com/actions/gh-actions-cache/service"
	"github.com/actions/gh-actions-cache/types"
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

func PrettyPrintCacheList(caches []types.ActionsCache) {
	for _, cache := range caches {
		var formattedRow string = getFormattedCacheInfo(cache)
		fmt.Println(formattedRow)
	}

}
func PrettyPrintTrimmedCacheList(caches []types.ActionsCache) {
	PrettyPrintCacheList(caches)
	if len(caches) > 30 {
		fmt.Printf("... and %d more\n\n", len(caches))
	}
	fmt.Print("\n")
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

func lastAccessedTime(lastAccessedAt string) string {
	var now time.Time = time.Now()
	lastAccessedTs, _ := time.Parse(time.RFC3339, lastAccessedAt)

	diff := now.Sub(lastAccessedTs)

	var lastAccessedTimeStr string = ""
	if diff.Minutes() < 1 {
		lastAccessedTimeStr = "Used less than a minute ago"
	} else if diff.Minutes() == 1 {
		lastAccessedTimeStr = "Used 1 minute ago"
	} else if diff.Minutes() < 60 {
		lastAccessedTimeStr = fmt.Sprintf("Used %d minutes ago", int(diff.Minutes()))
	} else if diff.Minutes() < 120 {
		lastAccessedTimeStr = "Used 1 hour ago"
	} else if diff.Hours() < 24 {
		lastAccessedTimeStr = fmt.Sprintf("Used %d hours ago", int(diff.Hours()))
	} else if diff.Hours() >= 24 && diff.Hours() < 48 {
		lastAccessedTimeStr = "Used 1 day ago"
	} else if diff.Hours() >= 48 {
		lastAccessedTimeStr = fmt.Sprintf("Used %d days ago", int(diff.Hours()/24))
	}
	return lastAccessedTimeStr
}

func getFormattedCacheInfo(cache types.ActionsCache) string {
	cacheKey := trimCacheKeyBasedOnWindowSize(cache.Key)
	return fmt.Sprintf(" %s\t [%s]\t %s\t %s", cacheKey, FormatCacheSize(cache.SizeInBytes), cache.Ref, lastAccessedTime(cache.LastAccessedAt))
}

func ListAllCaches(queryParams url.Values, key string, artifactCache service.ArtifactCacheService) []types.ActionsCache {
	listApiResponse := artifactCache.ListCaches(queryParams)
	caches := listApiResponse.ActionsCaches
	totalCaches := listApiResponse.TotalCount
	if totalCaches > 100 {
		for page := 2; page <= int(math.Ceil(float64(listApiResponse.TotalCount)/100)); page++ {
			queryParams.Set("page", strconv.Itoa(page))
			listApiResponse = artifactCache.ListCaches(queryParams)
			caches = append(caches, listApiResponse.ActionsCaches...)
		}
	}
	return caches
}

func RedTick() string {
	src := "\u2713"
	tick, _ := utf8.DecodeRuneInString(src)
	redTick := color.Colorize(color.Red, string(tick))
	return redTick
}
