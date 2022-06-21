package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var limit int

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("repo", "R", "", "Select another repository for finding actions cache.")
	listCmd.Flags().StringP("branch", "B", "", "Filter by branch")
	listCmd.Flags().IntVarP(&limit, "limit", "", 30, "Maximum number of items to fetch (default is 30, max limit is 100)")
	listCmd.Flags().StringP("key", "", "", "Filter by key")
	listCmd.Flags().StringP("order", "", "", "Order of caches returned (asc/desc)")
	listCmd.Flags().StringP("sort", "", "", "Sort fetched caches (used/size/created)")
	listCmd.SetHelpTemplate(getListHelp())
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists the actions cache",
	Long:  `Lists the actions cache`,
	Run: func(cmd *cobra.Command, args []string) {
		COMMAND = "list"
		r, _ := cmd.Flags().GetString("repo")
		branch, _ := cmd.Flags().GetString("branch")
		key, _ := cmd.Flags().GetString("key")
		order, _ := cmd.Flags().GetString("order")
		sort, _ := cmd.Flags().GetString("sort")

		repo, err := getRepo(r)
		if err != nil {
			log.Fatal(err)
		}

		totalCacheSize := getCacheUsage(repo)
		fmt.Printf("Total caches size %s\n", formatCacheSize(totalCacheSize))

		queryParams := generateQueryParams(branch, limit, key, order, sort)
		caches := listCaches(repo, queryParams)

		fmt.Printf("Showing %d of %d cache entries in %s/%s\n", totalShownCacheEntry(len(caches)), len(caches), repo.Owner(), repo.Name())
		for _, cache := range caches {
			fmt.Printf("%s\t [%s]\t %s\t %s\n", cache.Key, formatCacheSize(cache.Size), cache.Ref, cache.LastAccessedAt)
		}
	},
}

func totalShownCacheEntry(totalCaches int) int {
	if totalCaches < limit {
		return totalCaches
	}
	return limit
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
	--sort <string>				Sort fetched caches (used/size/created)

INHERITED FLAGS
	--help		Show help for command

EXAMPLES:
	$ gh actions-cache list
	$ gh actions-cache list --limit 100
	$ gh actions-cache list --order desc
`
}
