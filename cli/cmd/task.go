package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	taskCmd = &cobra.Command{
		Use: "task",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				fmt.Println("Failed to show help:", err)
				os.Exit(1)
			}
		},
	}
	listCmd = &cobra.Command{
		Use: "list",
		Run: doList,
	}
	addCmd = &cobra.Command{
		Use: "add",
		Run: doAdd,
	}
	viewCmd = &cobra.Command{
		Use: "view",
		Run: doView,
	}
	updateCmd = &cobra.Command{
		Use: "update",
		Run: doUpdate,
	}
	deleteCmd = &cobra.Command{
		Use: "delete",
		Run: doDelete,
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

func doList(cmd *cobra.Command, args []string) {

}

func doAdd(cmd *cobra.Command, args []string) {

}

func doView(cmd *cobra.Command, args []string) {

}

func doUpdate(cmd *cobra.Command, args []string) {

}

func doDelete(cmd *cobra.Command, args []string) {

}
