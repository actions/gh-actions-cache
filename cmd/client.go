package cmd

import (
	"fmt"
	"log"
	"net/url"

	gh "github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	ghRepo "github.com/cli/go-gh/pkg/repository"
)

type cacheInfo struct {
	Key            string
	Ref            string
	LastAccessedAt string
	Size           float64
}

func getCacheUsage(repo ghRepo.Repository) float64 {
	client, err := getRestClient(repo)
	if err != nil {
		log.Fatal(err)
	}
	pathComponent := fmt.Sprintf("repos/%s/%s/actions/cache/usage", repo.Owner(), repo.Name())
	var apiResults map[string]interface{}
	err = client.Get(pathComponent, &apiResults)
	if err != nil {
		log.Fatal(err)
	}

	cacheSizeResult := apiResults["active_caches_size_in_bytes"].(float64)
	return cacheSizeResult
}

func listCaches(repo ghRepo.Repository, queryParams url.Values) []cacheInfo {
	client, err := getRestClient(repo)
	if err != nil {
		log.Fatal(err)
	}
	pathComponent := fmt.Sprintf("repos/%s/%s/actions/caches", repo.Owner(), repo.Name())
	var apiResults map[string]interface{}
	err = client.Get(pathComponent+"?"+queryParams.Encode(), &apiResults)
	if err != nil {
		log.Fatal(err)
	}

	actionsCachesResult := apiResults["actions_caches"].([]interface{})

	var caches []cacheInfo
	for _, item := range actionsCachesResult {
		caches = append(caches, cacheInfo{
			Key:            item.(map[string]interface{})["key"].(string),
			Ref:            item.(map[string]interface{})["ref"].(string),
			LastAccessedAt: item.(map[string]interface{})["last_accessed_at"].(string),
			Size:           item.(map[string]interface{})["size_in_bytes"].(float64),
		})
	}
	return caches
}

func deleteCaches(repo ghRepo.Repository, queryParams url.Values) float64 {
	client, err := getRestClient(repo)
	if err != nil {
		log.Fatal(err)
	}
	pathComponent := fmt.Sprintf("repos/%s/%s/actions/caches", repo.Owner(), repo.Name())
	var apiResults map[string]interface{}
	err = client.Delete(pathComponent+"?"+queryParams.Encode(), &apiResults)
	if err != nil {
		if err.(api.HTTPError).StatusCode == 404 {
			return 0
		} else {
			log.Fatal(err)
		}
	}

	totalDeletedCachesResult := apiResults["total_count"].(float64)
	return totalDeletedCachesResult
}

func getRestClient(repo ghRepo.Repository) (api.RESTClient, error) {
	opts := api.ClientOptions{
		Host:    repo.Host(),
		Headers: map[string]string{"User-Agent": fmt.Sprintf("gh-actions-cache/%s/%s", VERSION, COMMAND)},
	}
	client, err := gh.RESTClient(&opts)
	if err != nil {
		return nil, err
	}
	return client, nil
}
