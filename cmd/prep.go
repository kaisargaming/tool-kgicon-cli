/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"strings"

	"os"

	"github.com/spf13/cobra"
)

// prepCmd represents the prep command
var prepCmd = &cobra.Command{
	Use:   "prep",
	Short: "Prepare icon files",
	Long: `Varied on the provider, prep command will prepare
svg files for kg/icons package, using upstreams submodule as
source files.`,
	Args: cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		provider := args[0]
		if err := execute(provider); err != nil {
			return err
		}
		return nil
	},
}

func execute(provider string) error {
	if provider == "hero" {
		return heroCopy()
	} else if provider == "majestic" {
		return majesticCopy()
	}

	return nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func majesticCopy() error {
	var err error
	var fds []os.FileInfo

	spath := "./resources/upstreams/majesticons"
	dpath := "./resources/providers/majestic"

	if err := checkDir(spath); !err {
		return errors.New("Upstream submodule hero is not available in resources/upstreams/majesticons")
	}

	if err := checkDir(dpath + "/o"); !err {
		var create error
		if create = os.MkdirAll(dpath+"/o", 0755); create != nil {
			return create
		}
		if create = os.MkdirAll(dpath+"/s", 0755); create != nil {
			return create
		}
		fmt.Println("Destination dir created " + dpath)
	}

	fmt.Println("Copying files from " + spath + " to " + dpath)

	if err = Dir(spath+"/solid", dpath+"/s"); err != nil {
		return err
	}

	if err = Dir(spath+"/line", dpath+"/o"); err != nil {
		return err
	}

	// rename
	if fds, err = ioutil.ReadDir(dpath + "/o"); err != nil {
		return err
	}

	for _, fd := range fds {
		srcName := fd.Name()
		split := strings.Split(srcName, "-")
		last := split[len(split)-1]

		if last == "line.svg" {
			dstName := strings.Join(split[:len(split)-1], "-") + ".svg"
			if rename := os.Rename(dpath+"/o/"+srcName, dpath+"/o/"+dstName); rename != nil {
				return rename
			}
			fmt.Println("Renaming " + srcName + " to " + dstName)
		}
	}

	return nil
}

func heroCopy() error {
	var err error

	spath := "./resources/upstreams/heroicons/optimized"
	dpath := "./resources/providers/hero"

	spathExist := checkDir(spath)
	dpathExist := checkDir(dpath)

	if !spathExist {
		return errors.New("Upstream submodule hero is not available in resources/upstreams/heroicons")
	}

	if !dpathExist {
		// create dir if not exist
		if err = os.MkdirAll(dpath, 0755); err != nil {
			return err
		}
		fmt.Println("Destination dir created " + dpath)
	}

	fmt.Println("Copying files from " + spath + " to " + dpath)

	if err = Dir(spath, dpath); err != nil {
		return err
	}

	if err = os.Rename(dpath+"/outline", dpath+"/o"); err != nil {
		return err
	}

	if err = os.Rename(dpath+"/solid", dpath+"/s"); err != nil {
		return err
	}

	return nil
}

// Dir copies a whole directory recursively
func Dir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = Dir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = File(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

// File copies a single file from src to dst
func File(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func checkDir(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

func init() {
	rootCmd.AddCommand(prepCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// prepCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// prepCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
