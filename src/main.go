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
		Use:   "clink",
		Short: "Welcome to CLink Synchronizer, Version : 1.2.1",
	}
	var syncCmd = &cobra.Command{
		Use:   "sync",
		Short: "Kafka log data to Clickhouse storage connection program, its performance is excellent and easy to use",
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
	syncCmd.CompletionOptions.HiddenDefaultCmd = true
	syncCmd.Flags().BoolP("debug", "d", false, "Enable debugging")
	syncCmd.Flags().StringP("config", "c", "./config.yaml", "Configuration file path")
	rootCmd.AddCommand(syncCmd)

	var regularCmd = &cobra.Command{
		Use:   "reg",
		Short: "Regex Test",
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
	regularCmd.Flags().StringP("regular", "r", "", "Regular expression")
	regularCmd.Flags().StringP("file", "f", "./test.log", "Debugging log files")
	rootCmd.AddCommand(regularCmd)

	err := rootCmd.Execute()
	if err != nil {
		color.Red(err.Error())
	}
}
