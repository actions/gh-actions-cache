package service

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
	ListCaches(repo ghRepo.Repository, queryParams url.Values) types.ListApiResponse
    DeleteCaches(repo ghRepo.Repository, queryParams url.Values) int
}

type ArtifactCache struct{
    HttpClient	api.RESTClient
}

func NewArtifactCache(repo ghRepo.Repository, command string, version string) ArtifactCacheService {
	opts := api.ClientOptions{
		Host:    repo.Host(),
		Headers: map[string]string{"User-Agent": fmt.Sprintf("gh-actions-cache/%s/%s", version, command)},
	}
	restClient, _ := gh.RESTClient(&opts)
	return &ArtifactCache{restClient}
}

func (a *ArtifactCache) GetCacheUsage(repo ghRepo.Repository) float64 {
	pathComponent := fmt.Sprintf("repos/%s/%s/actions/cache/usage", repo.Owner(), repo.Name())
	var apiResults types.RepoLevelUsageApiResponse
	err := a.HttpClient.Get(pathComponent, &apiResults)
	if err != nil {
		log.Fatal(err)
	}

	return apiResults.ActiveCacheSizeInBytes
}

func (a *ArtifactCache) ListCaches(repo ghRepo.Repository, queryParams url.Values) types.ListApiResponse {
	pathComponent := fmt.Sprintf("repos/%s/%s/actions/caches", repo.Owner(), repo.Name())
	var apiResults types.ListApiResponse
	err := a.HttpClient.Get(pathComponent+"?"+queryParams.Encode(), &apiResults)
	if err != nil {
		log.Fatal(err)
	}

	return apiResults
}

func (a *ArtifactCache) DeleteCaches(repo ghRepo.Repository, queryParams url.Values) int {
	pathComponent := fmt.Sprintf("repos/%s/%s/actions/caches", repo.Owner(), repo.Name())
	var apiResults types.DeleteApiResponse
	err := a.HttpClient.Delete(pathComponent+"?"+queryParams.Encode(), &apiResults)
	if err != nil {
		log.Fatal(err)
	}

	return apiResults.TotalCount
}