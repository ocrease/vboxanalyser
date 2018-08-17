package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ocrease/vboxanalyser/pkg/server"
	"github.com/ocrease/vboxanalyser/pkg/vbo"
	"github.com/pkg/browser"
)

const VboxExtension = ".vbo"

func main() {
	cli := flag.Bool("cli", false, "Enable command line interface only")
	dir := flag.String("dir", ".", "Specify the directory to scan")
	channel := flag.String("c", "LOT_Engine_Spd", "Specify the channel to analyse - rpm, speedKph, speedMph")
	threshold := flag.Float64("t", 8300, "Specify the RPM threshold")

	flag.Parse()

	if !*cli {
		browser.OpenURL("http://localhost:8080")
		server.NewServer().Start()
	}

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
			return err
		}
		if !info.IsDir() {
			if filepath.Ext(path) == VboxExtension {
				file := vbo.ParseFile(path)
				//fmt.Printf("%v - num points %v, num columns %v\n", path, len(file.Data.Rows), len(file.Columns))
				v, err := file.MaxValue(channel)
				if err != nil {
					fmt.Printf("%v - %v\n", path, err)
				}
				if v > threshold {
					fmt.Printf("%v - %v laps - %v\n", path, vbo.NumLaps(&file), v)
				}
			}
		}
		return nil
	}
}
