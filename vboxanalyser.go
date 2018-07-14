package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ocrease/vboxanalyser/models"
)

const VboxExtension = ".vbo"

func main() {
	dir := flag.String("dir", ".", "Specify the directory to scan")
	channel := flag.String("c", "rpm", "Specify the channel to analyse - rpm, speedKph, speedMph")
	threshold := flag.Float64("t", 8300, "Specify the RPM threshold")

	flag.Parse()

	path, _ := filepath.Abs(*dir)
	fmt.Printf("Analysing .vbo files in: %v\n", path)

	err := filepath.Walk(*dir, createFileProcessor(*channel, *threshold))
	if err != nil {
		log.Fatal(err)
	}

}

func createFileProcessor(channel string, threshold float64) func(string, os.FileInfo, error) error {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			if filepath.Ext(path) == VboxExtension {
				file := models.ParseFile(path)
				if v := file.MaxValueWithFunc(models.ExtractValueFunctionFactory(channel)); v > threshold {
					fmt.Printf("%v - %v\n", path, v)
				}
			}
		}
		return nil
	}
}
