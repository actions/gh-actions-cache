package cmd

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/actions/gh-actions-cache/internal"
	"github.com/actions/gh-actions-cache/service"
	"github.com/actions/gh-actions-cache/types"
	"github.com/spf13/cobra"
)

func NewCmdDelete() *cobra.Command {
	COMMAND = "delete"
	f := types.InputFlags{}

	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete cache by key",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Printf("accepts 1 arg(s), received %d\n", len(args))
				return
			}
			key := args[0]

			repo, err := internal.GetRepo(f.Repo)
			if err != nil {
				fmt.Println(err)
				return
			}
			artifactCache := service.NewArtifactCache(repo, COMMAND, VERSION)
			queryParams := internal.GenerateQueryParams(f.Branch, 100, key, "", "", 1)

			if !f.Confirm {
				var matchedCaches = getCacheListWithExactMatch(queryParams, key, artifactCache)
				if len(matchedCaches) == 0 {
					fmt.Printf("Cache with input key '%s' does not exist\n", key)
					return
				}
				fmt.Printf("You're going to delete %d cache ", len(matchedCaches))
				if len(matchedCaches) == 1 {
					fmt.Printf("entry\n\n")
				} else {
					fmt.Printf("entries\n\n")
				}
				internal.PrettyPrintTrimmedCacheList(matchedCaches)
				choice := ""
				prompt := &survey.Select{
					Message: "Are you sure you want to delete the cache entries?",
					Options: []string{"Delete", "Cancel"},
				}
				err := survey.AskOne(prompt, &choice)
				if err != nil {
					fmt.Println("Error occured while taking input from user while trying to delete cache")
					return
				}
				f.Confirm = choice == "Delete"
				fmt.Println()
			}
			if f.Confirm {
				var sb strings.Builder
				cachesDeleted := artifactCache.DeleteCaches(queryParams)
				if cachesDeleted > 0 {
					sb.WriteString(internal.RedTick())
					sb.WriteString(" Deleted ")
					sb.WriteString(strconv.Itoa(cachesDeleted))
					sb.WriteString(" cache ")
					if cachesDeleted == 1 {
						sb.WriteString("entry")
					} else {
						sb.WriteString("entries")
					}
					sb.WriteString(" with key ")
					sb.WriteString(key)
				} else {
					sb.WriteString("Cache with input key '")
					sb.WriteString(key)
					sb.WriteString("' does not exist")
				}
				fmt.Println(sb.String())
				sb.Reset()
			}
		},
	}
	deleteCmd.Flags().StringVarP(&f.Repo, "repo", "R", "", "Select another repository for finding actions cache.")
	deleteCmd.Flags().StringVarP(&f.Branch, "branch", "B", "", "Filter by branch")
	deleteCmd.Flags().BoolVar(&f.Confirm, "confirm", false, "Delete the cache without asking user for confirmation.")
	deleteCmd.SetHelpTemplate(getDeleteHelp())

	return deleteCmd
}

func getDeleteHelp() string {
	return `
gh-actions-cache: Works with GitHub Actions Cache. 

USAGE:
	gh actions-cache delete <key> [flags]

ARGUMENTS:
	key		cache key which needs to be deleted
	
FLAGS:
	-R, --repo <[HOST/]owner/repo>		Select another repository using the [HOST/]OWNER/REPO format
	-B, --branch <string>			Filter by branch
	--confirm				Confirm deletion without prompting

INHERITED FLAGS
	--help		Show help for command

EXAMPLES:
	$ gh actions-cache delete Linux-node-f5dbf39c9d11eba80242ac13
`
}

func getCacheListWithExactMatch(queryParams url.Values, key string, artifactCache service.ArtifactCacheService) []types.ActionsCache {
	caches := internal.ListAllCaches(queryParams, key, artifactCache)
	var exactMatchedKeys []types.ActionsCache
	for _, cache := range caches {
		if strings.EqualFold(key, cache.Key) {
			exactMatchedKeys = append(exactMatchedKeys, cache)
		}
	}
	return exactMatchedKeys
}
