package service

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/actions/gh-actions-cache/types"
	"gopkg.in/h2non/gock.v1"
	"github.com/actions/gh-actions-cache/internal"
)

const VERSION string = "0.0.1"

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
	assert.Nil(t, err)

	artifactCache := NewArtifactCache(repo, "list", VERSION)
	totalCacheSize, err := artifactCache.GetCacheUsage()

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
	assert.Nil(t, err)

	artifactCache := NewArtifactCache(repo, "list", VERSION)
	totalCacheSize, err := artifactCache.GetCacheUsage()

	assert.NotNil(t, err)
	assert.Equal(t, totalCacheSize ,float64(-1))
	assert.True(t, gock.IsDone(), printPendingMocks(gock.Pending()))
}

func TestListCaches_Success(t *testing.T) {
	t.Cleanup(gock.Off)

	gock.New("https://api.github.com").
		Get("/repos/testOrg/testRepo/actions/caches").
		Reply(200).
		JSON(`{
			"total_count": 1,
			"actions_caches": [
				{
					"id": 29,
					"ref": "refs/heads/master",
					"key": "Linux-build-cache-node-modules-3fd22dd3a926d576e2562e8b76a5ff157cd3b986f3d44195acfe7efa6bc05919-8",
					"version": "7fcda33c1e1d849a13bcc06f49b9ab64efc01ca9dabe4d7a8d0d387feef4fc88",
					"last_accessed_at": "2022-06-22T20:32:45.550000000Z",
					"created_at": "2022-06-22T20:32:45.550000000Z",
					"size_in_bytes": 2432967
				}]
			}`)


	repo, err := internal.GetRepo("testOrg/testRepo")
	assert.Nil(t, err)

	artifactCache := NewArtifactCache(repo, "list", VERSION)
	queryParams := internal.GenerateQueryParams("", 30, "", "", "", 1)
	listCacheResponse, err := artifactCache.ListCaches(queryParams)

	assert.Nil(t, err)
	assert.NotNil(t, listCacheResponse)
	assert.Equal(t, listCacheResponse.TotalCount ,1)
	assert.Equal(t, len(listCacheResponse.ActionsCaches) ,1)
	assert.Equal(t, listCacheResponse.ActionsCaches[0].Id ,29)
	assert.True(t, gock.IsDone(), printPendingMocks(gock.Pending()))
}

func TestListCaches_Failure(t *testing.T) {
	t.Cleanup(gock.Off)

	gock.New("https://api.github.com").
		Get("/repos/testOrg/testRepo/actions/caches").
		Reply(404).
		JSON(`{
			"message": "Not Found",
			"documentation_url": "https://docs.github.com/rest/reference/actions#get-github-actions-cache-list-for-a-repository"
		}`)


	repo, err := internal.GetRepo("testOrg/testRepo")
	assert.Nil(t, err)

	artifactCache := NewArtifactCache(repo, "list", VERSION)
	queryParams := internal.GenerateQueryParams("", 30, "", "", "", 1)
	listCacheResponse, err := artifactCache.ListCaches(queryParams)

	assert.NotNil(t, err)
	assert.Equal(t, listCacheResponse ,types.ListApiResponse{})
	assert.True(t, gock.IsDone(), printPendingMocks(gock.Pending()))
}

func TestDeleteCaches_Success(t *testing.T) {
	t.Cleanup(gock.Off)

	gock.New("https://api.github.com").
		Delete("/repos/testOrg/testRepo/actions/caches").
		Reply(200).
		JSON(`{
			"total_count": 1,
			"actions_caches": [
				{
					"id": 29,
					"ref": "refs/heads/master",
					"key": "Linux-build-cache-node-modules-3fd22dd3a926d576e2562e8b76a5ff157cd3b986f3d44195acfe7efa6bc05919-8",
					"version": "7fcda33c1e1d849a13bcc06f49b9ab64efc01ca9dabe4d7a8d0d387feef4fc88",
					"last_accessed_at": "2022-06-22T20:32:45.550000000Z",
					"created_at": "2022-06-22T20:32:45.550000000Z",
					"size_in_bytes": 2432967
				}]
			}`)


	repo, err := internal.GetRepo("testOrg/testRepo")
	assert.Nil(t, err)

	artifactCache := NewArtifactCache(repo, "delete", VERSION)
	queryParams := internal.GenerateQueryParams("", 30, "", "", "", 1)
	deletedCache, err := artifactCache.DeleteCaches(queryParams)

	assert.Nil(t, err)
	assert.Equal(t, deletedCache ,1)
	assert.True(t, gock.IsDone(), printPendingMocks(gock.Pending()))
}

func TestDeleteCaches_Failure(t *testing.T) {
	t.Cleanup(gock.Off)

	gock.New("https://api.github.com").
		Delete("/repos/testOrg/testRepo/actions/caches").
		Reply(404).
		JSON(`{
			"message": "Not Found",
			"documentation_url": "https://docs.github.com/rest/reference/actions#get-github-actions-cache-list-for-a-repository"
		}`)


	repo, err := internal.GetRepo("testOrg/testRepo")
	assert.Nil(t, err)

	artifactCache := NewArtifactCache(repo, "delete", VERSION)
	queryParams := internal.GenerateQueryParams("", 30, "", "", "", 1)
	deletedCache, err := artifactCache.DeleteCaches(queryParams)

	assert.NotNil(t, err)
	assert.Equal(t, deletedCache ,-1)
	assert.True(t, gock.IsDone(), printPendingMocks(gock.Pending()))
}

func printPendingMocks(mocks []gock.Mock) string {
	paths := []string{}
	for _, mock := range mocks {
		paths = append(paths, mock.Request().URLStruct.String())
	}
	return fmt.Sprintf("%d unmatched mocks: %s", len(paths), strings.Join(paths, ", "))
}