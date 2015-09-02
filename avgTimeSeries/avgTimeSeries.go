// avgTimeseries expects a list of files in two column format. The script then averages
// items in column 1 on a per row basis whereas column 0 is assumed to be identical
// between all files (in effect the time of events counted in row 2).
package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: avgTimeseries <directory|file names>")
		os.Exit(1)
	}

	fileNames := collectFiles()
	if len(fileNames) == 0 {
		fmt.Fprintf(os.Stderr, "no files provided\n")
		os.Exit(1)
	}

	avg := readCurrent(fileNames[0])
	for _, f := range fileNames[1:] {
		val := readCurrent(f)
		if len(val) != len(avg) {
			fmt.Fprintf(os.Stderr, "incorrect file length for %s\n", f)
		}
		for i := range val {
			avg[i] += val[i]
		}
	}

	// normalize
	numItems := float64(len(fileNames))
	for i := range avg {
		avg[i] /= numItems
		fmt.Println(avg[i])
	}
}

// readCurrent reads a two column data file and returns the second column as
// an array of float64's
func readCurrent(fileName string) []float64 {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
	defer file.Close()

	var reader io.Reader = file
	if filepath.Ext(fileName) == ".gz" {
		reader, err = gzip.NewReader(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}
	scanner := bufio.NewScanner(reader)

	var nums []float64
	prev := 0.0
	for scanner.Scan() {
		items := strings.Fields(scanner.Text())
		if len(items) != 2 {
			fmt.Fprintf(os.Stderr, "incorrect number of columns")
			os.Exit(1)
		}
		num, err := strconv.ParseFloat(items[1], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error converting string %s into float", items[1])
			os.Exit(1)
		}
		nums = append(nums, -(num - prev))
		prev = num
	}
	return nums
}

func collectFiles() []string {
	var fileList []string
	for _, fileName := range os.Args[1:] {

		info, err := os.Lstat(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			continue
		}

		if info.IsDir() {
			file, err := os.Open(fileName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
			content, err := file.Readdir(-1)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
			for _, item := range content {
				pathName := path.Join(fileName, item.Name())
				fileList = append(fileList, pathName)
			}
		} else {
			fileList = append(fileList, fileName)
		}
	}
	return fileList
}
