/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// createlistCmd represents the createlist command
var createlistCmd = &cobra.Command{
	Use:   "createlist setname path/to/dir",
	Short: "Create list of icon name",
	Long: `Creates list of icon names for each icon sets
and save them in a json file which will be published to
public/kgicons/.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var apperr error
		var abspath string

		name := args[0]
		dir := args[1]

		if abspath, apperr = filepath.Abs(dir); apperr != nil {
			return apperr
		}

		if _, apperr = os.Stat(abspath); os.IsNotExist(apperr) {
			return apperr
		}
		if apperr = create(name, abspath); apperr != nil {
			return apperr
		}
		return nil
	},
}

func create(name, abspath string) error {
	var err error
	var files []string
	var names []string

	//fmt.Println(abspath)

	err = filepath.Walk(abspath,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() == false {
				files = append(files, path)
			}
			return nil
		})

	if err != nil {
		return err
	}

	for _, file := range files {
		fileSplit := strings.Split(file, "/")
		prefix := fileSplit[len(fileSplit)-2]
		svgname := strings.Split(fileSplit[len(fileSplit)-1], ".")
		names = append(names, prefix+":"+svgname[0])
	}

	// fmt.Println(names)

	jsondata, _ := json.Marshal(names)

	//fmt.Println(string(data))

	jsonfname := "resources/js/" + name + ".json"
	if err = ioutil.WriteFile(jsonfname, jsondata, 0755); err != nil {
		return err
	}

	fmt.Println("File " + jsonfname + " Created")

	return nil
}

func init() {
	rootCmd.AddCommand(createlistCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createlistCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createlistCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
