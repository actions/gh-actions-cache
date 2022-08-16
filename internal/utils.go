package internal

import (
	"errors"
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/TwiN/go-color"
	"github.com/actions/gh-actions-cache/types"
	gh "github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	ghRepo "github.com/cli/go-gh/pkg/repository"
	ghTableprinter "github.com/cli/go-gh/pkg/tableprinter"
	ghTerm "github.com/cli/go-gh/pkg/term"
	"github.com/nleeper/goment"
)

const MB_IN_BYTES = 1024 * 1024
const GB_IN_BYTES = 1024 * 1024 * 1024

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
	terminal := ghTerm.FromEnv()
	w, _, _ := terminal.Size()
	tp := ghTableprinter.New(terminal.Out(), terminal.IsTerminalOutput(), w)

	for _, cache := range caches {
		tp.AddField(cache.Key)
		tp.AddField(FormatCacheSize(cache.SizeInBytes))
		tp.AddField(cache.Ref)
		tp.AddField(lastAccessedTime(cache.LastAccessedAt))
		tp.EndRow()
	}

	_ = tp.Render()
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
	return lastAccessed.FromNow()
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

func HttpErrorHandler(err error, errMsg404 string) types.HandledError {
	var httpError api.HTTPError
	if errors.As(err, &httpError) && httpError.StatusCode == 404 {
		return types.HandledError{Message: errMsg404, InnerError: err}
	} else if errors.As(err, &httpError) && httpError.StatusCode >= 400 && httpError.StatusCode < 500 {
		return types.HandledError{Message: httpError.Message, InnerError: err}
	} else {
		return types.HandledError{Message: "We could not process your request due to internal error.", InnerError: err}
	}
}
