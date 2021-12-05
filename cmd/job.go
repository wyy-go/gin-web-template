package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/wyy-go/go-web-template/internal/jobs"
)

var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "run job",
	PreRun: func(cmd *cobra.Command, args []string) {
		log.Println("PreRun")
		setup()
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("do nothing")
	},
}

var jobOnceCmd = &cobra.Command{
	Use:   "once",
	Short: "run job once",
	Long:  `run job once.`,
	Args:  cobra.MinimumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		log.Println("PreRun")
		setup()
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if job, ok := jobs.Jobs[name]; ok {
			job.Job()
		}
	},
}

var jobListCmd = &cobra.Command{
	Use:   "list",
	Short: "list all jobs",
	Long:  `list all jobs.`,
	Run: func(cmd *cobra.Command, args []string) {
		for k, v := range jobs.Jobs {
			fmt.Printf("%s [%s]\n", k, v.Spec)
		}
	},
}

func init() {
	jobCmd.AddCommand(
		jobListCmd,
		jobOnceCmd,
	)

	rootCmd.AddCommand(jobCmd)
}
