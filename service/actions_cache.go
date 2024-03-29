package service

import (
	"fmt"
	"math"
	"net/url"
	"strconv"

	"github.com/actions/gh-actions-cache/types"
	gh "github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	ghRepo "github.com/cli/go-gh/pkg/repository"
)

type ArtifactCacheService interface {
	GetCacheUsage() (float64, error)
	ListCaches(queryParams url.Values) (types.ListApiResponse, error)
	DeleteCaches(queryParams url.Values) (int, error)
	ListAllCaches(queryParams url.Values, key string) ([]types.ActionsCache, error)
}

type ArtifactCache struct {
	HttpClient api.RESTClient
	repo       ghRepo.Repository
}

func NewArtifactCache(repo ghRepo.Repository, command string, version string) (ArtifactCacheService, error) {
	opts := api.ClientOptions{
		Host:    repo.Host(),
		Headers: map[string]string{"User-Agent": fmt.Sprintf("gh-actions-cache/%s/%s", version, command)},
	}
	restClient, err := gh.RESTClient(&opts)
	if err != nil {
		return nil, err
	}
	return &ArtifactCache{HttpClient: restClient, repo: repo}, nil
}

func (a *ArtifactCache) GetCacheUsage() (float64, error) {
	pathComponent := fmt.Sprintf("repos/%s/%s/actions/cache/usage", a.repo.Owner(), a.repo.Name())
	var apiResults types.RepoLevelUsageApiResponse
	err := a.HttpClient.Get(pathComponent, &apiResults)
	if err != nil {
		return -1, err
	}

	return apiResults.ActiveCacheSizeInBytes, nil
}

func (a *ArtifactCache) ListCaches(queryParams url.Values) (types.ListApiResponse, error) {
	pathComponent := fmt.Sprintf("repos/%s/%s/actions/caches", a.repo.Owner(), a.repo.Name())
	var apiResults types.ListApiResponse
	err := a.HttpClient.Get(pathComponent+"?"+queryParams.Encode(), &apiResults)

	if err != nil {
		return types.ListApiResponse{}, err
	}

	return apiResults, nil
}

func (a *ArtifactCache) DeleteCaches(queryParams url.Values) (int, error) {
	pathComponent := fmt.Sprintf("repos/%s/%s/actions/caches", a.repo.Owner(), a.repo.Name())
	var apiResults types.DeleteApiResponse
	err := a.HttpClient.Delete(pathComponent+"?"+queryParams.Encode(), &apiResults)
	if err != nil {
		return 0, err
	}
	return apiResults.TotalCount, nil
}

func (a *ArtifactCache) ListAllCaches(queryParams url.Values, key string) ([]types.ActionsCache, error) {
	var listApiResponse types.ListApiResponse
	listApiResponse, err := a.ListCaches(queryParams)
	if err != nil {
		return nil, err
	}

	caches := listApiResponse.ActionsCaches
	totalCaches := listApiResponse.TotalCount
	if totalCaches > 100 {
		for page := 2; page <= int(math.Ceil(float64(listApiResponse.TotalCount)/100)); page++ {
			queryParams.Set("page", strconv.Itoa(page))
			listApiResponse, err := a.ListCaches(queryParams)
			if err != nil {
				return nil, err
			}
			caches = append(caches, listApiResponse.ActionsCaches...)
		}
	}
	return caches, nil
}
