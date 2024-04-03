package main

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/xxscloud5722/cLink/src/app"
	"os"
)

func main() {
	var rootCmd = &cobra.Command{
		Use: "clink",
		//Long: "Kafka log data to Clickhouse storage connection program, its performance is excellent and easy to use",
		Run: func(cmd *cobra.Command, args []string) {
			configPath, err := cmd.Flags().GetString("config")
			if err != nil {
				color.Red(err.Error())
			}
			debug, err := cmd.Flags().GetBool("debug")
			if err != nil {
				color.Red(err.Error())
			}
			if _, err = os.Stat(configPath); os.IsNotExist(err) {
				err = cmd.Help()
				if err != nil {
					color.Red(err.Error())
				}
				return
			}
			err = app.Run(&configPath, debug)
			if err != nil {
				color.Red(err.Error())
			}
		},
	}
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debugging")
	rootCmd.PersistentFlags().StringP("config", "c", "./config.yaml", "Configuration file path")
	color.Blue("Welcome to CLink Synchronizer, Version : 1.2.1")
	color.Blue("CLink, created by Cat, is a robust solution for synchronizing data from Kafka messaging systems to ClickHouse databases.")
	err := rootCmd.Execute()
	if err != nil {
		color.Red(err.Error())
	}
}
