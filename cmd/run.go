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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	valid, err := validateConfigSchema()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	if !valid {
		log.Print("config JSON did not pass the schema test")
		os.Exit(1)
	}

	fmt.Println("run called")

	if len(args) == 0 {
		fmt.Println("Command needed to run. Please run 'list' to see available commands.")
		return
	}

	configMap, err := readConfigFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Config map", configMap)

	// TODO: Less hard-coding
	commandMap := configMap.Commands

	log.Println("Command map: ", commandMap)

	fmt.Println("Args", args)

	commandToRun := args[0]

	fmt.Println("Attempting to run: ", commandToRun)

	commandInfo, exists := commandMap[commandToRun]
	if !exists {
		fmt.Println(commandToRun, " is not an available command. Please run 'list' to see available commands.")
		return
	}

	fmt.Println("Command exists and has the definition of: ", commandMap)

	runCommands := commandInfo.RunCommands
	cmdArgs := []string{}
	if runCommands != nil {
		cmdArgs = runCommands
	}

	workingdir, _ := os.Getwd()
	commandLocation := filepath.Join(workingdir, commandInfo.Command)
	fmt.Println("Command location: ", commandLocation)
	cmdArgs = append(cmdArgs, commandLocation)
	cmdArgs = append(cmdArgs, args[1:]...)
	cmdVar := exec.Command(commandInfo.Lang, cmdArgs...)
	cmdVar.Env = os.Environ()

	if cmdOut, err = cmdVar.Output(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	result := string(cmdOut)

	fmt.Println("Command output: ", result)
}

// Config struct mapping the fp config json
type Config struct {
	Commands commands `json:"commands"`
	// TODO: ADD OTHER CONFIG PROPERTIES HERE AS WELL
}

type commands map[string]Command

// Command struct represents info relating to each command
type Command struct {
	Command        string   `json:"command"`
	Environment    string   `json:"environment"`
	Lang           string   `json:"lang"`
	RunCommands    []string `json:"runCommands"`
	InstallCommand string   `json:"install"`
}

func validateConfigSchema() (valid bool, err error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return false, err
	}

	// https://github.com/xeipuuv/gojsonschema/issues/92
	// reference must have a leader reference (e.g. 'file://')
	schemaLocation := fmt.Sprintf("file://%s/config-schema.json", workingDir)
	log.Printf("Loading Config schema from %s", schemaLocation)

	// https://github.com/xeipuuv/gojsonschema
	// TODO: We can load these configs from a lot of different places
	// TODO: but we are doing locally for the time being
	schemaLoader := gojsonschema.NewReferenceLoader(schemaLocation)

	configLocation := fmt.Sprintf("file://%s/config.json", workingDir)
	log.Printf("Loading config from %s", configLocation)
	configLoader := gojsonschema.NewReferenceLoader(configLocation)

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

func readConfigFile() (Config, error) {
	jsonFile, err := os.Open("config.json")

	if err != nil {
		fmt.Println("Could not read config file (config.json)", err)
		return Config{}, err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config Config
	json.Unmarshal([]byte(byteValue), &config)

	return config, nil
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
