package actionsCacheClient

import (
	"fmt"
	"log"
	"net/url"

	"github.com/cli/go-gh/pkg/api"
	ghRepo "github.com/cli/go-gh/pkg/repository"
)

type cacheInfo struct {
	Key            string
	Ref            string
	LastAccessedAt string
	Size           float64
}

type ListApiResponse struct {
	TotalCount    int            `json:"total_count"`
	ActionsCaches []ActionsCache `json:"actions_caches"`
}

type ActionsCache struct {
	Id             int      `json:"id"`
	Ref            string   `json:"ref"`
	Key            string   `json:"key"`
	Version        string   `json:"version"`
	LastAccessedAt string	`json:"last_accessed_at"`
	CreatedAt      string 	`json:"created_at"`
	SizeInBytes    float64  `json:"size_in_bytes"`
}

func GetCacheUsage(repo ghRepo.Repository, client api.RESTClient) float64 {
	pathComponent := fmt.Sprintf("repos/%s/%s/actions/cache/usage", repo.Owner(), repo.Name())
	var apiResults map[string]interface{}
	err := client.Get(pathComponent, &apiResults)
	if err != nil {
		log.Fatal(err)
	}

	cacheSizeResult := apiResults["active_caches_size_in_bytes"].(float64)
	return cacheSizeResult
}

func ListCaches(repo ghRepo.Repository, queryParams url.Values, client api.RESTClient) []cacheInfo {
	pathComponent := fmt.Sprintf("repos/%s/%s/actions/caches", repo.Owner(), repo.Name())
	var apiResults ListApiResponse
	err := client.Get(pathComponent+"?"+queryParams.Encode(), &apiResults)
	if err != nil {
		log.Fatal(err)
	}

	actionsCachesResult := apiResults.ActionsCaches

	var caches []cacheInfo
	for _, item := range actionsCachesResult {
		fmt.Println(item)
		caches = append(caches, cacheInfo{
			Key:            item.Key,
			Ref:            item.Ref,
			LastAccessedAt: item.LastAccessedAt,
			Size:           item.SizeInBytes,
		})
	}
	return caches
}

func DeleteCaches(repo ghRepo.Repository, queryParams url.Values, client api.RESTClient) float64 {
	pathComponent := fmt.Sprintf("repos/%s/%s/actions/caches", repo.Owner(), repo.Name())
	var apiResults map[string]interface{}
	err := client.Delete(pathComponent+"?"+queryParams.Encode(), &apiResults)
	if err != nil {
		log.Fatal(err)
	}

	totalDeletedCachesResult := apiResults["total_count"].(float64)
	return totalDeletedCachesResult
}
