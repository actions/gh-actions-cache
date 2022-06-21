package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/actions/gh-actions-cache/internal"
	"github.com/actions/gh-actions-cache/actionsCacheClient"
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("repo", "R", "", "Select another repository for finding actions cache.")
	listCmd.Flags().StringP("branch", "B", "", "Filter by branch")
	listCmd.Flags().IntP("limit", "", 30, "Maximum number of items to fetch (default is 30, max limit is 100)")
	listCmd.Flags().StringP("key", "", "", "Filter by key")
	listCmd.Flags().StringP("order", "", "", "Order of caches returned (asc/desc)")
	listCmd.Flags().StringP("sort", "", "", "Sort fetched caches (last-used/size/created-at)")
	listCmd.SetHelpTemplate(getListHelp())
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists the actions cache",
	Long:  `Lists the actions cache`,
	Run: func(cmd *cobra.Command, args []string) {
		COMMAND = "list"

		if len(args) != 0 {
			fmt.Printf("Invalid argument(s). Expected 0 received %d\n", len(args))
			fmt.Println(getListHelp())
			return
		}

		r, _ := cmd.Flags().GetString("repo")
		branch, _ := cmd.Flags().GetString("branch")
		limit, _ := cmd.Flags().GetInt("limit")
		key, _ := cmd.Flags().GetString("key")
		order, _ := cmd.Flags().GetString("order")
		sort, _ := cmd.Flags().GetString("sort")

		repo, err := internal.GetRepo(r)
		if err != nil {
			log.Fatal(err)
		}

		validateInputs(sort, order, limit)

		client, err := internal.GetRestClient(repo, "0.0.1", "list")
		if err != nil {
			log.Fatal(err)
		}

		if branch == "" && key == "" {
			totalCacheSize := actionsCacheClient.GetCacheUsage(repo, client)
			fmt.Printf("Total caches size %s\n\n", internal.FormatCacheSize(totalCacheSize))
		}

		queryParams := internal.GenerateQueryParams(branch, limit, key, order, sort)
		caches := actionsCacheClient.ListCaches(repo, queryParams, client)

		fmt.Printf("Showing %d of %d cache entries in %s/%s\n\n", totalShownCacheEntry(len(caches), limit), len(caches), repo.Owner(), repo.Name())
		for _, cache := range caches {
			fmt.Printf("%s\t [%s]\t %s\t %s\n", cache.Key, internal.FormatCacheSize(cache.Size), cache.Ref, cache.LastAccessedAt)
		}
	},
}

func totalShownCacheEntry(totalCaches int, limit int) int {
	if totalCaches < limit {
		return totalCaches
	}
	return limit
}

func validateInputs(sort string, order string, limit int){
	if order != "" && order != "asc" && order != "desc"{
		log.Fatal(fmt.Errorf(fmt.Sprintf("%s is not a valid value for order flag. Allowed values: asc/desc", order)))
	}

	if sort != "" && sort != "last-used" && sort != "size" && sort != "created-at"{
		log.Fatal(fmt.Errorf(fmt.Sprintf("%s is not a valid value for sort flag. Allowed values: last-used/size/created-at", sort)))
	}

	if limit < 1{
		log.Fatal(fmt.Errorf(fmt.Sprintf("%d is not a valid value for limit flag. Allowed values: > 1", limit)))
	}
}

func getListHelp() string {
	return `
gh-actions-cache: Works with GitHub Actions Cache. 

USAGE:
	gh actions-cache list [flags]

ARGUMENTS:
	No Arguments

FLAGS:
	-R, --repo <[HOST/]owner/repo>		Select another repository using the [HOST/]OWNER/REPO format
	-B, --branch <string>			Filter by branch
	-L, --limit <int>			Maximum number of items to fetch (default is 30, max limit is 100)
	--key <string>				Filter by key
	--order <string>			Order of caches returned (asc/desc)
	--sort <string>				Sort fetched caches (last-used/size/created-at)

INHERITED FLAGS
	--help		Show help for command

EXAMPLES:
	$ gh actions-cache list
	$ gh actions-cache list --limit 100
	$ gh actions-cache list --order desc
`
}
