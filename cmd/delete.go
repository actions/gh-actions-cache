package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringP("repo", "R", "", "Select another repository for finding actions cache.")
	deleteCmd.Flags().StringP("branch", "B", "", "Filter by branch")
	deleteCmd.SetHelpTemplate(getDeleteHelp())
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		branch, _ := cmd.Flags().GetString("branch")
		fmt.Println("DELETE")
		fmt.Println(repo)
		fmt.Println(branch)
	},
}

func getDeleteHelp() string {
	return `
gh-actions-cache: Works with GitHub Actions Cache. 

USAGE:
	gh actions-cache lisy [flags]

ARGUMENTS:
	key		cache key which needs to be deleted
	
FLAGS:
	-R, --repo <[HOST/]owner/repo>		Select another repository
	-B, --branch <string>			Filter by branch

INHERITED FLAGS
	--help		Show help for command

EXAMPLES:
	$ gh actions-cache delete Linux-node-f5dbf39c9d11eba80242ac13
`
}
