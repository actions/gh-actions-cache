package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("repo", "R", "", "Select another repository for finding actions cache.")
	listCmd.Flags().StringP("branch", "B", "", "Filter by branch")
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
		repo, _ := cmd.Flags().GetString("repo")
		branch, _ := cmd.Flags().GetString("branch")
		key, _ := cmd.Flags().GetString("key")
		order, _ := cmd.Flags().GetString("order")
		sort, _ := cmd.Flags().GetString("sort")
		fmt.Println("LIST")
		fmt.Println(repo)
		fmt.Println(branch)
		fmt.Println(key)
		fmt.Println(order)
		fmt.Println(sort)
	},
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
