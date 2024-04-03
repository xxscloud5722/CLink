package main

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/xxscloud5722/cLink/src/app"
	"os"
	"regexp"
	"strconv"
	"strings"
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

	var regularCmd = &cobra.Command{
		Use: "reg",
		Run: func(cmd *cobra.Command, args []string) {
			regular, err := cmd.Flags().GetString("regular")
			if err != nil {
				color.Red(err.Error())
			}
			file, err := cmd.Flags().GetString("file")
			if err != nil {
				color.Red(err.Error())
			}
			// 文件如果不存在
			if _, err = os.Stat(file); os.IsNotExist(err) {
				return
			}
			// 读取文件
			textBytes, err := os.ReadFile(file)
			if err != nil {
				return
			}
			text := string(textBytes)
			if regular == "" || text == "" {
				return
			}
			color.Magenta("Regexp: " + regular)
			color.Magenta("Text: " + text)
			matches := regexp.MustCompile(regular).FindStringSubmatch(text)
			if matches == nil {
				return
			}
			color.Magenta("======================================================")
			for i, item := range matches {
				if i <= 0 {
					continue
				}
				color.Blue("[" + strconv.Itoa(i) + "]" + strings.Trim(item, "\n"))
			}
		},
	}
	regularCmd.PersistentFlags().StringP("regular", "r", "", "Regular expression")
	regularCmd.PersistentFlags().StringP("file", "f", "./test.log", "Debugging log files")
	rootCmd.AddCommand(regularCmd)

	color.Blue("Welcome to CLink Synchronizer, Version : 1.2.1")
	err := rootCmd.Execute()
	if err != nil {
		color.Red(err.Error())
	}
}
