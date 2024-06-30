package cmd

import (
	"cine-tool/app"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cine",
	Short: "Cine",
	Long:  `Cine Tool`,
	Run: func(cmd *cobra.Command, args []string) {
		app.RunServer()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println("Run Server Error: ", err)
		os.Exit(1)
	}
}
