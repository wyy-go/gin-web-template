package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use: "go-web-template",
	Short:             "go-web-template server",
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

func exitError(msg interface{}) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "conf/dev.yaml", "config file")
}

func Execute() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		_ = rootCmd.Help()
	}

	if err := rootCmd.Execute(); err != nil {
		exitError(err)
	}
}
