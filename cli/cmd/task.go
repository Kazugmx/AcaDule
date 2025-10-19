package cmd

import (
	"acadule-cli/internal/acaduleapi"
	"acadule-cli/internal/config"
	"acadule-cli/internal/simpleform"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
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
	cfg := loadConfig()

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
	cfg := loadConfig()

	fmt.Println("Fetching status...")
	title := simpleform.Ask("What's your new task name?")
	requestData := acaduleapi.TaskAddRequest{
		Title: title,
	}

	res, err := acaduleapi.Add(apiURL, cfg.Token, requestData)
	if err != nil {
		fmt.Println("Failed to request api:", err)
		os.Exit(1)
	}

	fmt.Println("Task created!")
	fmt.Println("- Title:", title)
	fmt.Println("- Id:", res.Id)
	fmt.Println("- Status:", res.Status)
}

func doView(cmd *cobra.Command, args []string) {
	cfg := loadConfig()

	if len(args) == 0 {
		fmt.Println("Set id as first argument.")
		os.Exit(1)
	}

	id := args[0]
	data, err := acaduleapi.View(apiURL, cfg.Token, id)
	if err != nil {
		fmt.Println("Failed to get task data:", err)
		os.Exit(1)
	}

	fmt.Println("---- Task Information ----")
	fmt.Println("Title:", data.Title)
	fmt.Println("Status:", data.Progress)
}

func doUpdate(cmd *cobra.Command, args []string) {
	cfg := loadConfig()

	if len(args) == 0 {
		fmt.Println("Set id as first argument.")
		os.Exit(1)
	}
	id := args[0]

	requestData := acaduleapi.UpdateRequest{
		Id: id,
	}
	var selectedFieldIds []string
	huh.NewMultiSelect[string]().
		Title("Which fields want to edit?").
		Value(&selectedFieldIds).
		Options(
			huh.NewOption("Title", "title"),
			huh.NewOption("Description", "description"),
			huh.NewOption("Progress", "progress"),
			huh.NewOption("Deadline", "deadline"),
			huh.NewOption("HasDone", "hasDone"),
		).Run()

	for _, fieldId := range selectedFieldIds {
		switch fieldId {
		case "title":
			requestData.Title = simpleform.Ask("New title")
		case "description":
			requestData.Description = simpleform.Ask("New description")
		case "progress":
			huh.NewSelect[acaduleapi.TaskProgress]().
				Title("New task progress").
				Value(&requestData.Progress).
				Options(
					huh.NewOption("Not Started", acaduleapi.NOT_STARTED),
					huh.NewOption("In Progress", acaduleapi.IN_PROGRESS),
					huh.NewOption("Complete", acaduleapi.COMPLETE),
					huh.NewOption("Suspended", acaduleapi.SUSPENDED),
				).Run()
		case "deadline":
			fmt.Println("Sorry. Deadline can't update from terminal at this time.")
		case "hasDone":
			requestData.HasDone = simpleform.Confirm("Task has done?")
		}
	}
	res, err := acaduleapi.Update(apiURL, cfg.Token, requestData)
	if err != nil {
		fmt.Println("Failed to request update api:", err)
		os.Exit(1)
	}

	fmt.Println("Updated! id:", res.Id)
}

func doDelete(cmd *cobra.Command, args []string) {
	cfg := loadConfig()

	if len(args) == 0 {
		fmt.Println("Set id as first argument.")
		os.Exit(1)
	}
	id := args[0]

	requestData := acaduleapi.DeleteRequest{
		Id: id,
	}

	if !simpleform.Confirm("Are you really remove the task?") {
		fmt.Println("See you!")
		os.Exit(0)
	}

	err := acaduleapi.Delete(apiURL, cfg.Token, requestData)
	if err != nil {
		fmt.Println("Failed to delete task:", err)
		os.Exit(1)
	}

	fmt.Println("Successfully to delete task. id:", id)
}

func loadConfig() config.Config {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Error occurred on loading config:", err)
		os.Exit(1)
	}
	validateAndUpdateConfig(&cfg)
	return cfg
}
