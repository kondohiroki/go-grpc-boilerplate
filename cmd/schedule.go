package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/kondohiroki/go-grpc-boilerplate/config"
	"github.com/kondohiroki/go-grpc-boilerplate/internal/scheduler"
	"github.com/lnquy/cron"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddGroup(&cobra.Group{ID: "schedule", Title: "Schedule:"})
	rootCmd.AddCommand(
		listScheduleCommand,
		startScheduleCommand,
	)
}

var scheduleCommands map[string]*cobra.Command
var startScheduleCommand = &cobra.Command{
	Use:     "schedule:run",
	Short:   "Start schedule job",
	GroupID: "schedule",
	PreRun: func(cmd *cobra.Command, _ []string) {
		setUpConfig()
		setUpLogger()
		setUpPostgres()
		setUpRedis()
	},
	Run: func(cmd *cobra.Command, _ []string) {
		printScheduleList()
		printScheduleList()

		scheduleCommands = map[string]*cobra.Command{}

		scheduler.Start(cmd, scheduleCommands)
	},
}

var listScheduleCommand = &cobra.Command{
	Use:     "schedule:list",
	Short:   "List all schedule jobs",
	GroupID: "schedule",
	PreRun: func(cmd *cobra.Command, _ []string) {
		setUpConfig()
		setUpLogger()
	},
	Run: func(_ *cobra.Command, _ []string) {
		printScheduleList()
	},
}

func printScheduleList() {
	exprDesc, _ := cron.NewDescriptor()

	// Print the job list as a table in the console
	tableWriter := table.NewWriter()
	tableWriter.SetOutputMirror(os.Stdout)
	tableWriter.AppendHeader(table.Row{"No.", "Job Name", "Cron Expression", "Schedule"})
	for i, schedule := range config.GetConfig().Schedules {
		desc, _ := exprDesc.ToDescription(schedule.Cron, cron.Locale_en)

		tableWriter.AppendRow(table.Row{
			i + 1,
			schedule.Job,
			schedule.Cron,
			desc,
		})

	}

	tableWriter.Render()
}
