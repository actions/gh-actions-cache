package cmd

import (
	"fmt"
	"log"

	"github.com/actions/gh-actions-cache/internal"
	"github.com/actions/gh-actions-cache/service"
	"github.com/actions/gh-actions-cache/types"
	"github.com/spf13/cobra"
)

func NewCmdList() *cobra.Command {
	COMMAND = "list"

	f := types.InputFlags{}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists the actions cache",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 0 {
				fmt.Printf("Invalid argument(s). Expected 0 received %d\n", len(args))
				fmt.Println(getListHelp())
				return
			}

			repo, err := internal.GetRepo(f.Repo)
			if err != nil {
				log.Fatal(err)
			}

			validateInputs(f)

			artifactCache := service.NewArtifactCache(repo, COMMAND, VERSION)

			if f.Branch == "" && f.Key == "" {
				totalCacheSize := artifactCache.GetCacheUsage()
				fmt.Printf("Total caches size %s\n\n", internal.FormatCacheSize(totalCacheSize))
			}

			queryParams := internal.GenerateQueryParams(f.Branch, f.Limit, f.Key, f.Order, f.Sort, 1)
			listCacheResponse := artifactCache.ListCaches(queryParams)

			totalCaches := listCacheResponse.TotalCount
			caches := listCacheResponse.ActionsCaches

			fmt.Printf("Showing %d of %d cache entries in %s/%s\n\n", displayedEntriesCount(len(caches), f.Limit), totalCaches, repo.Owner(), repo.Name())
			for _, cache := range caches {
				fmt.Printf("%s\t [%s]\t %s\t %s\n", cache.Key, internal.FormatCacheSize(cache.SizeInBytes), cache.Ref, cache.LastAccessedAt)
			}
		},
	}

	listCmd.Flags().StringVarP(&f.Repo, "repo", "R", "", "Select another repository for finding actions cache.")
	listCmd.Flags().StringVarP(&f.Branch, "branch", "B", "", "Filter by branch")
	listCmd.Flags().IntVarP(&f.Limit, "limit", "", 30, "Maximum number of items to fetch (default is 30, max limit is 100)")
	listCmd.Flags().StringVarP(&f.Key, "key", "", "", "Filter by key")
	listCmd.Flags().StringVarP(&f.Order, "order", "", "", "Order of caches returned (asc/desc)")
	listCmd.Flags().StringVarP(&f.Sort, "sort", "", "", "Sort fetched caches (last-used/size/created-at)")
	listCmd.SetHelpTemplate(getListHelp())

	return listCmd
}

func displayedEntriesCount(totalCaches int, limit int) int {
	if totalCaches < limit {
		return totalCaches
	}
	return limit
}

func validateInputs(input types.InputFlags) {
	if input.Order != "" && input.Order != "asc" && input.Order != "desc" {
		log.Fatal(fmt.Errorf(fmt.Sprintf("%s is not a valid value for order flag. Allowed values: asc/desc", input.Order)))
	}

	if input.Sort != "" && input.Sort != "last-used" && input.Sort != "size" && input.Sort != "created-at" {
		log.Fatal(fmt.Errorf(fmt.Sprintf("%s is not a valid value for sort flag. Allowed values: last-used/size/created-at", input.Sort)))
	}

	if input.Limit < 1 || input.Limit > 100 {
		log.Fatal(fmt.Errorf(fmt.Sprintf("%d is not a valid value for limit flag. Allowed values: 1-100", input.Limit)))
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
