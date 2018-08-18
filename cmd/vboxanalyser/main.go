package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/ocrease/vboxanalyser/pkg/server"
	"github.com/ocrease/vboxanalyser/pkg/vbo"
	"github.com/pkg/browser"
)

var (
	cli       bool
	open      bool
	dir       string
	channel   string
	threshold float64
	port      int
)

func init() {
	flag.BoolVar(&cli, "cli", false, "Enable command line interface only")

	flag.IntVar(&port, "port", 8080, "HTTP Port - not in CLI mode")
	flag.BoolVar(&open, "open", true, "Open Browser - not in CLI mode")
	flag.StringVar(&dir, "dir", ".", "Specify the directory to scan - CLI Mode only")
	flag.StringVar(&channel, "c", "rpm", "Specify the channel to analyse - rpm, speedKph - CLI Mode only")
	flag.Float64Var(&threshold, "t", 8300, "Specify the threshold - CLI Mode only")
	flag.Parse()
}

func main() {

	if !cli {
		if open {
			browser.OpenURL(fmt.Sprintf("http://localhost:%v", port))
		}
		server.NewServer(port).Start()
	}

	path, _ := filepath.Abs(dir)
	fmt.Printf("Analysing .vbo files in: %v\n", path)

	analyser := vbo.Analyser{}
	analyser.AnalyseDirectory(path, func(summary vbo.FileSummary) {
		var channelValue float64
		switch channel {
		case "rpm":
			channelValue = summary.MaxRpm
		case "speedKph":
			channelValue = summary.MaxVelocity
		}
		if channelValue > threshold {
			fmt.Printf("%v - %v laps - %v\n", summary.Path, summary.NumLaps, channelValue)
		}
	})
}
