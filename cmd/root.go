package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "wrpc service",
	Short: "wrpc service cli",
	Long:  "wrpc service cli",
}

// Execute cobra command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Unexpected execute error, err %v", err)
		os.Exit(1)
	}
}
