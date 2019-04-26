// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

type config struct {
	BasePath    string `yaml:"basePath"`
	CurrentPath string `yaml:"currentPath"`
}

// nbCmd represents the nb command
var nbCmd = &cobra.Command{
	Use:   "nb",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var c config
		c.getConf()
		fmt.Println(c.CurrentPath)

		name, _ := cmd.Flags().GetString("name")
		folderSwitch, _ := cmd.Flags().GetString("switch")

		//Creating new notebook/csv command
		if name != "" {
			path := filepath.Join(c.BasePath, "nb", name)
			basePath := c.BasePath
			createNotebook(basePath, path, name)
			createFile(name, path)

			fmt.Println(name + " folder created.")
		}

		//Switch to another notebook command
		if folderSwitch != "" {
			basePath := c.BasePath
			newPath := filepath.Join(c.BasePath, "nb", folderSwitch, folderSwitch+".csv")
			testYaml(basePath, newPath)
			fmt.Println("Switching to notebook - " + folderSwitch)
		}
	},
}

func createNotebook(basePath, path, name string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
	switchPath := filepath.Join(path, name+".csv")
	testYaml(basePath, switchPath)
}

func createFile(name, path string) {
	file, err := os.Create(filepath.Join(path, filepath.Base(name+".csv")))
	defer file.Close()

	if err != nil {
		os.Exit(1)
	}
}

func (c *config) getConf() *config {
	os.Chdir(".config")
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func testYaml(basePath, path string) {
	os.Chdir(".config")
	// configfolder, _ := os.Getwd()
	// configfolder = configfolder + ".config"
	//fmt.Println(configfolder)

	var err = os.Remove("config.yaml")
	if err != nil {
		fmt.Printf("[ERR] %v", err)
		return
	}

	file, err := os.Create("config.yaml")
	if err != nil {
		fmt.Printf("[ERR] %v", err)
		return
	}
	defer file.Close()

	//os.Chdir(".config")
	str := fmt.Sprintf("basePath: %s\ncurrentNotebook: %s", basePath, path)
	data := []byte(str)
	err = ioutil.WriteFile("config.yaml", data, 0644)
	if err != nil {
		fmt.Printf("[ERR %v", err)
	}
}

func init() {
	rootCmd.AddCommand(nbCmd)
	nbCmd.Flags().StringP("name", "n", "", "Set name of Folder")
	nbCmd.Flags().StringP("switch", "s", "", "Switch notebooks")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
