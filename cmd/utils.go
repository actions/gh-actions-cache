package cmd

import (
	"fmt"
	"math"
	"net/url"
	"os"
	"strconv"
	"time"

	gh "github.com/cli/go-gh"
	ghRepo "github.com/cli/go-gh/pkg/repository"
	"github.com/moby/term"
)

func parseInputFlags(branch string, limit int, key string, order string, sort string) url.Values {
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
		repo, err := ghRepo.Parse(r)
		if err != nil {
			return nil, err
		}
		return repo, nil
	}

	currentRepo, err := gh.CurrentRepository()
	if err != nil {
		return nil, err
	}
	return currentRepo, nil
}

func formatCacheSize(size_in_bytes float64) string {
	if size_in_bytes < 1024 {
		return fmt.Sprintf("%.2f Bytes", size_in_bytes)
	}

	if size_in_bytes < 1024*1024 {
		return fmt.Sprintf("%.2f KB", size_in_bytes/1024)
	}

	if size_in_bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", size_in_bytes/1024/1024)
	}

	return fmt.Sprintf("%.2f GB", size_in_bytes/1024/1024/1024)
}

func getCacheListWithExactMatch(repo ghRepo.Repository, queryParams url.Values) []cacheInfo {

	listApiResponse := listCaches(repo, queryParams)
	var exactMatchedKeys []cacheInfo
	for _, cache := range listApiResponse {
		if queryParams.Get("key") == cache.Key {
			exactMatchedKeys = append(exactMatchedKeys, cache)
		}
	}
	return exactMatchedKeys
}

func prettyPrintCacheList(caches []cacheInfo) {
	numberOfCaches := len(caches)
	fmt.Print("\nYou're going to delete " + strconv.Itoa(numberOfCaches) + " cache ")
	if numberOfCaches == 1 {
		fmt.Println("entry")
	} else {
		fmt.Println("entries")
	}
	fmt.Println()
	for _, cache := range caches {
		var cacheKey string = trimCacheKeyBasedOnWindowSize(cache.Key)

		var resultRow string = "  " + cacheKey + "     [" + formatCacheSize(cache.Size) + "]     " + cache.Ref[11:] + "     " + lastAccessedHour(cache.LastAccessedAt)
		fmt.Println(resultRow)
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
