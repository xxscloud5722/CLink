package app

//
//import (
//	"fmt"
//	"github.com/fatih/color"
//	"gopkg.in/yaml.v3"
//	"os"
//)
//
////func (config *Config) Columns() []string {
////	if config.columns == nil {
////		var columns []string
////		for _, it := range strings.Split(config.Fields, ",") {
////			columns = append(columns, strings.TrimSpace(it))
////		}
////		config.columns = columns
////	}
////	return config.columns
////}
////
////func (config *Config) Regexp() *regexp.Regexp {
////	if config.regexp == nil {
////		config.regexp = regexp.MustCompile(config.Pattern)
////	}
////	return config.regexp
////}
//
//var app Config
//
//// SetConfig GLOBAL variable
//func SetConfig(fileName string) {
//	yamlFile, err := os.ReadFile(fileName)
//	if err != nil {
//		fmt.Println("[Config] Error reading YAML file:", err)
//		return
//	}
//	var config Config
//	err = yaml.Unmarshal(yamlFile, &config)
//	if err != nil {
//		color.Red("[Config] Error parse YAML file: %v", err)
//		return
//	}
//	app = config
//}
//
//// Get GLOBAL variable
//func Get() *Config {
//	return &app
//}
