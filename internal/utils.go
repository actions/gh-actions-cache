package internal

import (
	"fmt"
	"math"
	"net/url"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/TwiN/go-color"
	"github.com/actions/gh-actions-cache/types"
	gh "github.com/cli/go-gh"
	ghRepo "github.com/cli/go-gh/pkg/repository"
	"github.com/moby/term"
	"github.com/nleeper/goment"
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
	fd := os.Stdin.Fd()
	ws, _ := term.GetWinsize(fd)
	width := math.Min(float64(ws.Width), 180)
	keyWidth := int(math.Max(math.Floor(0.30*width), 6))
	sizeWidth := int(math.Max(math.Floor(0.12*width), 3))
	refWidth := int(math.Max(math.Floor(0.20*width), 4))
	timeWidth := int(math.Max(math.Floor(0.20*width), 4))
	for _, cache := range caches {
		var formattedRow string = getFormattedCacheInfo(cache, keyWidth, sizeWidth, refWidth, timeWidth)
		fmt.Println(formattedRow)
	}
}
func PrettyPrintTrimmedCacheList(caches []types.ActionsCache) {
	length := len(caches)
	limit := 30
	if length > limit {
		PrettyPrintCacheList(caches[:limit])
		fmt.Printf("... and %d more\n\n", length-limit)
	} else {
		PrettyPrintCacheList(caches[:length])
	}
	fmt.Print("\n")
}

func lastAccessedTime(lastAccessedAt string) string {
	lastAccessed, _ := goment.New(lastAccessedAt)
	return fmt.Sprintf("Used %s", lastAccessed.FromNow())
}

func trimOrPad(value string, maxSize int) string {
	if len(value) > maxSize {
		value = value[:maxSize-3] + "..."
	} else {
		value = value + strings.Repeat(" ", maxSize-len(value))
	}
	return value
}

func getFormattedCacheInfo(cache types.ActionsCache, keyWidth int, sizeWidth int, refWidth int, timeWidth int) string {
	key := trimOrPad(cache.Key, keyWidth)
	size := trimOrPad(fmt.Sprintf("[%s]", FormatCacheSize(cache.SizeInBytes)), sizeWidth)
	ref := trimOrPad(cache.Ref, refWidth)
	time := trimOrPad(lastAccessedTime(cache.LastAccessedAt), timeWidth)
	return fmt.Sprintf(" %s    %s    %s    %s", key, size, ref, time)
}

func RedTick() string {
	src := "\u2713"
	tick, _ := utf8.DecodeRuneInString(src)
	redTick := color.Colorize(color.Red, string(tick))
	return redTick
}

func PrintSingularOrPlural(count int, singularStr string, pluralStr string) string {
	if count == 1 {
		return fmt.Sprintf("%d %s", count, singularStr)
	}
	return fmt.Sprintf("%d %s", count, pluralStr)
}
