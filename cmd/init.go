/*
Copyright © 2023 Agustin Larreinegabe

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new GitHub organization configuration",
	Run: func(cmd *cobra.Command, args []string) {
		org, _ := cmd.Flags().GetString("org")
		dir, _ := cmd.Flags().GetString("dir")
		initConfigAAC(org, dir)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("org", "o", "", "Name of the GitHub organization")
	initCmd.Flags().StringP("dir", "d", "", "Working directory (default is the current directory)")
	initCmd.MarkFlagRequired("org")

	initCmd.MarkFlagFilename("dir", "yaml")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfigAAC(org string, dir string) {
	// If no directory is specified, use the current directory
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Create the organization directory within the specified or current directory
	dirName := filepath.Join(dir, org)
	err := os.Mkdir(dirName, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Change to the new directory
	err = os.Chdir(dirName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// // Initialize Git repo
	// cmd := exec.Command("git", "init")
	// err = cmd.Run()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// Create access-config.yaml
	file, err := os.Create("access-config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// (Optional) Write a basic structure to access-config.yaml
	file.WriteString("organizations:\n")
	fmt.Println("Initialization completed in", dirName)
}
