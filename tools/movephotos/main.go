package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/missdeer/golib/fsutil"

	flag "github.com/spf13/pflag"
)

var (
	sourceDir        string
	targetDir        string
	verbose          bool
	version          bool
	dryRun           bool
	removeSourceFile bool
)

func readDir(dirname string) error {
	fis, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Println(err)
		return err
	}
	for _, fi := range fis {
		sourcePath := filepath.Join(dirname, fi.Name())
		if fi.IsDir() {
			readDir(sourcePath)
		} else {
			targetPath := filepath.Join(targetDir,
				strconv.Itoa(fi.ModTime().Year()), fi.ModTime().Month().String(), strconv.Itoa(fi.ModTime().Day()),
				fi.Name())
			dir := path.Dir(targetPath)
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				if verbose {
					log.Println(dir, "not exists, create it now")
				}
				if !dryRun {
					os.MkdirAll(dir, 0644)
				}
			}
			if !dryRun {
				tfi, err := os.Stat(targetPath)
				if os.IsNotExist(err) || tfi.Size() != fi.Size() {
					os.Remove(targetPath)
					_, err := fsutil.CopyFile(sourcePath, targetPath)
					if err != nil {
						log.Println(err)
						continue
					}
				}
				os.Chtimes(targetPath, fi.ModTime(), fi.ModTime())
				if removeSourceFile {
					if err = os.Remove(sourcePath); err != nil {
						log.Println(err)
					}
				}
			}
			if verbose {
				fmt.Println(sourcePath, "move to", targetPath)
			}
		}
	}
	return nil
}

func main() {
	flag.StringVarP(&sourceDir, "source", "s", "", "set source directory path")
	flag.StringVarP(&targetDir, "target", "t", ".", "set target directory path")
	flag.BoolVarP(&verbose, "verbose", "V", false, "print verbose logs")
	flag.BoolVarP(&version, "version", "v", false, "print version information")
	flag.BoolVarP(&dryRun, "dry-run", "d", false, "dry run")
	flag.BoolVarP(&removeSourceFile, "remove-source-file", "r", false, "remove source file if copied file successfully")
	flag.Parse()

	if version {
		fmt.Println("Move photos")
		return
	}

	readDir(sourceDir)
}
