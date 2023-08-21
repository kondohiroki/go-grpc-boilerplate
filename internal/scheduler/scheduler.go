package scheduler

import (
	"fmt"
	"time"
	_ "time/tzdata"

	"github.com/go-co-op/gocron"
	"github.com/kondohiroki/go-grpc-boilerplate/config"
	"github.com/kondohiroki/go-grpc-boilerplate/internal/logger"
	"github.com/spf13/cobra"
)

var Timezone = time.Now().Location()

func Start(cmd *cobra.Command, commands map[string]*cobra.Command) {
	if config.GetConfig().Scheduler.Timezone != "" {
		Timezone, _ = time.LoadLocation(config.GetConfig().Scheduler.Timezone)
	}

	s := gocron.NewScheduler(Timezone)
	s.SingletonModeAll()

	for _, schedule := range config.GetConfig().Schedules {
		if schedule.IsEnabled {
			command, exists := commands[schedule.Job]
			if !exists {
				fmt.Printf("Command for job %s does not exist\n", schedule.Job)
				continue
			}

			err := runTask(s, cmd, schedule, command)
			if err != nil {
				fmt.Printf("Failed to schedule %s job: %v\n", schedule.Job, err)
				continue
			}

		} else {
			fmt.Printf("Job %s is not enabled\n", schedule.Job)
		}
	}

	fmt.Printf("Total jobs: %d jobs scheduled to run\n", len(s.Jobs()))
	fmt.Printf("Timezone: %s\n", s.Location().String())
	fmt.Println("Starting scheduler... (press Ctrl+C to quit)")

	s.StartImmediately()
	s.StartBlocking()
}

func runTask(s *gocron.Scheduler, cmd *cobra.Command, schedule config.Schedule, command *cobra.Command) error {
	task, err := s.CronWithSeconds(schedule.Cron).Do(getJobFunction(cmd, schedule, command))
	if err != nil {
		return err
	}

	// Set up event listeners
	task.SetEventListeners(func() {
		fmt.Printf("\n%s Job started -- round: %d\n", schedule.Job, task.RunCount())

	}, func() {
		time.Sleep(1 * time.Second)

		// Print next run time in both utc and asia/bangkok
		asiaBangkok, _ := time.LoadLocation("Asia/Bangkok")
		fmt.Printf("%s Next run: %s / %s\n", schedule.Job, task.NextRun().UTC().String(), task.NextRun().In(asiaBangkok).String())

	})
	return nil
}

func getJobFunction(cmd *cobra.Command, schedule config.Schedule, command *cobra.Command) func() {
	return func() {
		logger.Log.Info(fmt.Sprintf("Running %s job...", schedule.Job))
		command.Run(cmd, nil)
	}
}
