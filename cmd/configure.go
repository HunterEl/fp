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

	"bufio"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: configure,
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func configure(cmd *cobra.Command, args []string) {
	log.Print("Running the configure command...")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("[cmds repo]: ")
	repo, _ := reader.ReadString('\n')

	log.Printf("commands repo is %s", repo)
	// write the commands repo to fp.rc file
	// f, err := os.Create("../fp.rc")
	rcContents := fmt.Sprintf("commandsRepo=%s", repo)

	bytes := []byte(rcContents)
	err = ioutil.WriteFile("./fp.rc", bytes, 666)
	check(err)

	// defer f.Close()

	// w := bufio.NewWriter(f)
	// bytesWritten, err := w.WriteString(rcContents)
	// log.Printf("Wrote %d bytes", bytesWritten)
	// err = w.Flush()
	// check(err)
}

func init() {
	rootCmd.AddCommand(configureCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
