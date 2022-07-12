package internal

import (
	"fmt"
	"math"
	"os"
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
const SIZE_COLUMN_WIDTH = 15
const LAST_ACCESSED_AT_COLUMN_WIDTH = 20

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
	fd := os.Stdout.Fd()
	ws, _ := term.GetWinsize(fd)
	width := math.Max(math.Min(float64(ws.Width), 180), 100)

	sizeWidth := SIZE_COLUMN_WIDTH             // hard-coded size as the content is scoped
	timeWidth := LAST_ACCESSED_AT_COLUMN_WIDTH // hard-coded size as the content is scoped
	keyWidth := int(math.Floor(0.65 * float64(width-15-20)))
	refWidth := int(math.Floor(0.20 * float64(width-15-20)))

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
	return fmt.Sprintf(" %s", lastAccessed.FromNow())
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
	return fmt.Sprintf("%s %s %s %s", key, size, ref, time)
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
