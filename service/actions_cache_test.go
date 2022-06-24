package service

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"github.com/actions/gh-actions-cache/internal"
)

func TestGetCacheUsage_CorrectRepo(t *testing.T) {
	t.Cleanup(gock.Off)

	gock.New("https://api.github.com").
		Get("/repos/testOrg/testRepo/actions/cache/usage").
		Reply(200).
		JSON(`{
			"full_name": "t-dedah/vipul-bugbash",
			"active_caches_size_in_bytes": 291205,
			"active_caches_count": 12
		}`)


	repo, err := internal.GetRepo("testOrg/testRepo")
	if err != nil {
		log.Fatal(err)
	}

	artifactCache := NewArtifactCache(repo, "list", "0.0.1")
	totalCacheSize := artifactCache.GetCacheUsage()

	assert.Equal(t, totalCacheSize ,float64(291205))
	assert.Nil(t, err)
	assert.True(t, gock.IsDone(), printPendingMocks(gock.Pending()))
}

func TestGetCacheUsage_IncorrectRepo(t *testing.T) {
	t.Cleanup(gock.Off)

	gock.New("https://api.github.com").
		Get("/repos/testOrg/wrongRepo/actions/cache/usage").
		Reply(404).
		JSON(`{
			"message": "Not Found",
			"documentation_url": "https://docs.github.com/rest/reference/actions#get-github-actions-cache-usage-for-a-repository"
		}`)


	repo, err := internal.GetRepo("testOrg/wrongRepo")
	if err != nil {
		log.Fatal(err)
	}

	artifactCache := NewArtifactCache(repo, "list", "0.0.1")
	totalCacheSize := artifactCache.GetCacheUsage()

	fmt.Println(totalCacheSize)
	// assert.Equal(t, totalCacheSize ,float64(291205))
	assert.Nil(t, err)
	assert.True(t, gock.IsDone(), printPendingMocks(gock.Pending()))
}

func printPendingMocks(mocks []gock.Mock) string {
	paths := []string{}
	for _, mock := range mocks {
		paths = append(paths, mock.Request().URLStruct.String())
	}
	return fmt.Sprintf("%d unmatched mocks: %s", len(paths), strings.Join(paths, ", "))
}