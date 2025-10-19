package cmd

import "github.com/spf13/cobra"

var (
	taskCmd = &cobra.Command{
		Use: "task",
	}
	listCmd = &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	addCmd = &cobra.Command{
		Use: "add",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	viewCmd = &cobra.Command{
		Use: "view",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	updateCmd = &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	deleteCmd = &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
)

func init() {
	rootCmd.AddCommand(taskCmd)
	taskCmd.AddCommand(
		listCmd,
		addCmd,
		viewCmd,
		updateCmd,
		deleteCmd,
	)
}
