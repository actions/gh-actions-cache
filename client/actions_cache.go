package client

import (
	"fmt"
	"log"
	"net/url"

	"github.com/cli/go-gh/pkg/api"
	ghRepo "github.com/cli/go-gh/pkg/repository"
	gh "github.com/cli/go-gh"
	"github.com/actions/gh-actions-cache/types"
)

type ArtifactCacheService interface{
    GetCacheUsage(repo ghRepo.Repository) float64
	ListCaches(repo ghRepo.Repository, queryParams url.Values) []types.CacheInfo
    DeleteCaches(repo ghRepo.Repository, queryParams url.Values) int
}

type ArtifactCache struct{
    HttpClient	api.RESTClient
}

func NewArtifactCache(opts api.ClientOptions) ArtifactCache {
	if opts.Host == ""{
		return ArtifactCache{nil}
	}
	client, err := gh.RESTClient(&opts)
	if err != nil {
		log.Fatal(err)
	}
	return ArtifactCache{client}
}

func (a *ArtifactCache) GetCacheUsage(repo ghRepo.Repository) float64 {
	pathComponent := fmt.Sprintf("repos/%s/%s/actions/cache/usage", repo.Owner(), repo.Name())
	var apiResults types.RepoLevelUsageApiResponse
	err := a.HttpClient.Get(pathComponent, &apiResults)
	if err != nil {
		log.Fatal(err)
	}

	cacheSizeResult := apiResults.ActiveCacheSizeInBytes
	return cacheSizeResult
}

func (a *ArtifactCache) ListCaches(repo ghRepo.Repository, queryParams url.Values) []types.CacheInfo {
	pathComponent := fmt.Sprintf("repos/%s/%s/actions/caches", repo.Owner(), repo.Name())
	var apiResults types.ListApiResponse
	err := a.HttpClient.Get(pathComponent+"?"+queryParams.Encode(), &apiResults)
	if err != nil {
		log.Fatal(err)
	}

	actionsCachesResult := apiResults.ActionsCaches

	var caches []types.CacheInfo
	for _, item := range actionsCachesResult {
		caches = append(caches, types.CacheInfo{
			Key:            item.Key,
			Ref:            item.Ref,
			LastAccessedAt: item.LastAccessedAt,
			Size:           item.SizeInBytes,
		})
	}
	return caches
}

func (a *ArtifactCache) DeleteCaches(repo ghRepo.Repository, queryParams url.Values) int {
	pathComponent := fmt.Sprintf("repos/%s/%s/actions/caches", repo.Owner(), repo.Name())
	var apiResults types.DeleteApiResponse
	err := a.HttpClient.Delete(pathComponent+"?"+queryParams.Encode(), &apiResults)
	if err != nil {
		log.Fatal(err)
	}

	totalDeletedCachesResult := apiResults.TotalCount
	return totalDeletedCachesResult
}