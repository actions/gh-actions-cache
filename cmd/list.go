package cmd

import (
	"fmt"

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
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf(fmt.Sprintf("Invalid argument(s). Expected 0 received %d", len(args)))
			}

			repo, err := internal.GetRepo(f.Repo)
			if err != nil {
				return err
			}

			err = validateInputs(f)

			if err != nil {
				return err
			}

			artifactCache := service.NewArtifactCache(repo, COMMAND, VERSION)

			if f.branch == "" && f.key == "" {
				totalCacheSize, err := artifactCache.GetCacheUsage()
				if err != nil {
					return err
				}
				fmt.Printf("Total caches size %s\n\n", internal.FormatCacheSize(totalCacheSize))
			}

			queryParams := internal.GenerateQueryParams(f.branch, f.limit, f.key, f.order, f.sort, 1)
			listCacheResponse, err := artifactCache.ListCaches(queryParams)
			if err != nil {
				return err
			}

			totalCaches := listCacheResponse.TotalCount
			caches := listCacheResponse.ActionsCaches

			fmt.Printf("Showing %d of %d cache entries in %s/%s\n\n", displayedEntriesCount(len(caches), f.Limit), totalCaches, repo.Owner(), repo.Name())
			internal.PrettyPrintCacheList(caches)
      return nil
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

func validateInputs(input types.InputFlags) error {
	if input.order != "" && input.order != "asc" && input.order != "desc" {
		return fmt.Errorf(fmt.Sprintf("%s is not a valid value for order flag. Allowed values: asc/desc", input.order))
	}

	if input.sort != "" && input.sort != "last-used" && input.sort != "size" && input.sort != "created-at" {
		return fmt.Errorf(fmt.Sprintf("%s is not a valid value for sort flag. Allowed values: last-used/size/created-at", input.sort))
	}

	if input.limit < 1 || input.limit > 100 {
		return fmt.Errorf(fmt.Sprintf("%d is not a valid value for limit flag. Allowed values: 1-100", input.limit))
	}
	return nil
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
