package cmd

import (
	"acadule-cli/internal/acaduleapi"
	"acadule-cli/internal/config"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	taskCmd = &cobra.Command{
		Use:   "task",
		Short: "Edit tasks in AcaDule platform",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				fmt.Println("Failed to show help:", err)
				os.Exit(1)
			}
		},
	}
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Show task list",
		Run:   doList,
	}
	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add the task",
		Run:   doAdd,
	}
	viewCmd = &cobra.Command{
		Use:   "view",
		Short: "View detail of the task",
		Run:   doView,
	}
	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update the task data",
		Run:   doUpdate,
	}
	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete the task",
		Run:   doDelete,
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
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Error occurred on loading config:", err)
		os.Exit(1)
	}
	validateAndUpdateConfig(&cfg)

	fmt.Println("Fetching status...")

	allTasks, err := acaduleapi.GetAll(apiURL, cfg.Token)
	if err != nil {
		fmt.Println("Failed to get tasks:", err)
		os.Exit(1)
	}

	for _, d := range *allTasks {
		fmt.Printf("- %s: %s\n", d.Id, d.Title)
	}
	fmt.Printf("Count: %d\n", len(*allTasks))
}

func doAdd(cmd *cobra.Command, args []string) {

}

func doView(cmd *cobra.Command, args []string) {

}

func doUpdate(cmd *cobra.Command, args []string) {

}

func doDelete(cmd *cobra.Command, args []string) {

}
