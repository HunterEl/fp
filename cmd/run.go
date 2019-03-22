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
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: commandRun,
}

var (
	cmdOut []byte
	err    error
)

func commandRun(cmd *cobra.Command, args []string) {
	configRepo, err := readAppConfig()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	log.Printf("Commands Repo is %s", configRepo)

	exists, err := localRepoExists(configRepo)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	log.Printf("repo %s exists locally: %v", configRepo, exists)

	if !exists {
		log.Printf("Fetcing %s...", configRepo)
		_, err := fetchRepo(configRepo)
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}

	}

	repoHasConfig, err := localRepoHasConfig(configRepo)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	if repoHasConfig == false {
		log.Printf("Repo %s does not have a root config file (config.json)", configRepo)
		os.Exit(1)
	}

	log.Printf("repo %s has config file: %v", configRepo, repoHasConfig)

	configMap, err := loadRepoConfig(configRepo)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Print("Config Map: ", configMap)

	repoPath, err := fullRepoPath(configRepo)
	if err != nil {
		log.Print(err)
		return
	}

	workingDir, err := os.Getwd()
	if err != nil {
		log.Print(err)
		return
	}

	schemaLocation := fmt.Sprintf("%s/config-schema.json", workingDir)
	configLocation := fmt.Sprintf("%s/config.json", repoPath)

	valid, err := validateConfigSchema(schemaLocation, configLocation)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	if !valid {
		log.Print("config JSON did not pass the schema test")
		os.Exit(1)
	}

	if len(args) == 0 {
		fmt.Println("Command needed to run. Please run 'list' to see available commands.")
		return
	}

	commandMap := configMap.Commands

	log.Print("Args: ", args)

	commandToRun := args[0]

	log.Printf("Attempting to run: %s", commandToRun)

	commandInfo, exists := commandMap[commandToRun]
	if !exists {
		fmt.Println(commandToRun, " is not an available command. Please run 'list' to see available commands.")
		return
	}

	log.Printf("Command exists and has the definition of: %s", commandInfo)

	runCommands := commandInfo.RunCommands
	cmdArgs := []string{}
	if runCommands != nil {
		cmdArgs = runCommands
	}

	// Full path of the script in the caches directory
	commandLocation := filepath.Join(repoPath, commandInfo.Command)
	log.Printf("Command location: %s", commandLocation)

	cmdArgs = append(cmdArgs, commandLocation)
	cmdArgs = append(cmdArgs, args[1:]...)
	cmdVar := exec.Command(commandInfo.Lang, cmdArgs...)
	cmdVar.Env = os.Environ()

	if cmdOut, err = cmdVar.Output(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	result := string(cmdOut)

	log.Printf("Command output: %s", result)
}

func validateConfigSchema(schemaLocation string, configLocation string) (valid bool, err error) {
	// https://github.com/xeipuuv/gojsonschema/issues/92
	// reference must have a leader reference (e.g. 'file://')
	canonicalSchemaLocation := fmt.Sprintf("file://%s", schemaLocation)
	// log.Printf("Loading Config schema from %s", schemaLocation)

	// https://github.com/xeipuuv/gojsonschema
	// TODO: We can load these configs from a lot of different places
	// TODO: but we are doing locally for the time being
	schemaLoader := gojsonschema.NewReferenceLoader(canonicalSchemaLocation)

	canonicalConfigLocation := fmt.Sprintf("file://%s", configLocation)

	// log.Printf("Loading config from %s", configLocation)
	configLoader := gojsonschema.NewReferenceLoader(canonicalConfigLocation)

	result, err := gojsonschema.Validate(schemaLoader, configLoader)
	if err != nil {
		return false, err
	}

	if result.Valid() {
		return true, nil
	}

	// Go through all the schema errors and construct a new error message
	schemaErrors := []string{}
	for _, desc := range result.Errors() {
		formattedError := fmt.Sprintf("- %s", desc)
		schemaErrors = append(schemaErrors, formattedError)
	}

	schemaError := errors.New(strings.Join(schemaErrors, "\n"))
	return false, schemaError
}

func init() {
	rootCmd.AddCommand(runCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
