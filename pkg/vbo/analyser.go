package vbo

import (
	"fmt"
	"os"
	"path/filepath"
)

type Analyser struct{}

const VboxExtension = ".vbo"

func (a *Analyser) AnalyseDirectory(path string) []FileSummary {
	summaries := make([]FileSummary, 0)
	fmt.Printf("Analysing .vbo files in: %v\n", path)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {

			if filepath.Ext(path) == VboxExtension {
				fmt.Printf("Analysing %v\n", path)

				file := ParseFile(path)
				rpm, err := file.MaxValue("LOT_Engine_Spd")
				if err != nil {
					rpm = 0
				}
				vel, err := file.MaxValue("velocity")
				if err != nil {
					vel = 0
				}
				summaries = append(summaries, FileSummary{Path: path, NumLaps: NumLaps(&file), MaxVelocity: vel, MaxRpm: rpm})
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error processing directory %v - %v\n", path, err)
	}
	fmt.Printf("Analysing %v .vbo files in: %v\n", len(summaries), path)
	return summaries
}

// func createFileProcessor(summaries *[]FileSummary) func(string, os.FileInfo, error) error {
// 	return func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if !info.IsDir() {

// 			if filepath.Ext(path) == VboxExtension {
// 				fmt.Printf("Analysing %v\n", path)

// 				file := ParseFile(path)
// 				rpm, err := file.MaxValue("rpm")
// 				if err != nil {
// 					rpm = 0
// 				}
// 				vel, err := file.MaxValue("velocity")
// 				if err != nil {
// 					vel = 0
// 				}
// 				summaries = append(summaries, FileSummary{Path: path, NumLaps: NumLaps(&file), MaxVelocity: vel, MaxRpm: rpm})
// 			}
// 		}
// 		return nil
// 	}
// }

type FileSummary struct {
	Path        string  `json:"path"`
	NumLaps     uint32  `json:"numlaps"`
	MaxVelocity float64 `json:"maxvelocity"`
	MaxRpm      float64 `json:"maxrpm"`
}
