package cmd

import (
	"fmt"
	"net/url"

	"github.com/actions/gh-actions-cache/internal"
	"github.com/actions/gh-actions-cache/service"
	"github.com/actions/gh-actions-cache/types"
	"github.com/spf13/cobra"
)

func NewCmdList() *cobra.Command {
	COMMAND = "list"

	f := types.ListOptions{}

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

			// This will silence the usage (help) message as they are not needed for errors beyond this point
			cmd.SilenceUsage = true

			err = f.Validate()
			if err != nil {
				return err
			}

			artifactCache, err := service.NewArtifactCache(repo, COMMAND, VERSION)
			if err != nil {
				return types.HandledError{Message: err.Error(), InnerError: err}
			}

			if f.Branch == "" && f.Key == "" {
				totalCacheSize, err := artifactCache.GetCacheUsage()
				if err == nil && totalCacheSize > 0 {
					fmt.Printf("Total caches size %s\n\n", internal.FormatCacheSize(totalCacheSize))
				}
			}

			queryParams := url.Values{}
			f.GenerateQueryParams(queryParams)
			listCacheResponse, err := artifactCache.ListCaches(queryParams)
			if err != nil {
				return internal.HttpErrorHandler(err, "The given repo does not exist.", "We could not process your request due to internal error.")
			}

			totalCaches := listCacheResponse.TotalCount
			caches := listCacheResponse.ActionsCaches
			if len(caches) > 0 {
				fmt.Printf("Showing %d of %d cache entries in %s/%s\n\n", displayedEntriesCount(len(caches), f.Limit), totalCaches, repo.Owner(), repo.Name())
				internal.PrettyPrintCacheList(caches)
			} else if f.Key != "" {
				fmt.Printf("Cache with input key '%s' does not exist\n", f.Key)
			} else {
				fmt.Printf("There are no Actions caches currently present in this repo\n")
			}
			return nil
		},
	}

	listCmd.Flags().StringVarP(&f.Repo, "repo", "R", "", "Select another repository for finding actions cache.")
	listCmd.Flags().StringVarP(&f.Branch, "branch", "B", "", "Filter by branch")
	listCmd.Flags().IntVarP(&f.Limit, "limit", "", 30, "Number of items to fetch between 1 to 100")
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
