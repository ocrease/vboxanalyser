package vbo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Analyser struct{}

const VboxExtension = ".vbo"
const SummaryExtension = ".vbs"

func (a *Analyser) AnalyseDirectory(path string, consumer func(FileSummary)) {
	fmt.Printf("Analysing .vbo files in: %v\n", path)
	err := filepath.Walk(path, walkFunc(consumer))
	if err != nil {
		fmt.Printf("Error processing directory %v - %v\n", path, err)
	}
}

func walkFunc(consumer func(FileSummary)) func(path string, info os.FileInfo, err error) error {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {

			if filepath.Ext(path) == VboxExtension {
				handleFile(path, consumer)
			}
		}
		return nil
	}
}

func handleFile(path string, consumer func(FileSummary)) {
	fmt.Printf("Analysing %v\n", path)
	vbs, vbsExists := summaryExists(path)
	if vbsExists {
		if summary, err := loadSummary(vbs); err == nil {
			fmt.Println("Returning data from stored summary")
			consumer(*summary)
			return
		}
	}

	file := ParseFile(path)
	summary := generateSummary(&file)
	saveSummary(vbs, summary)
	consumer(*summary)
}

func generateSummary(file *File) *FileSummary {
	rpm, err := file.MaxValue("LOT_Engine_Spd")
	if err != nil {
		rpm = 0
	}
	vel, err := file.MaxValue("velocity")
	if err != nil {
		vel = 0
	}
	return &FileSummary{Path: file.Path, NumLaps: NumLaps(file), MaxVelocity: vel, MaxRpm: rpm}
}

func summaryExists(path string) (string, bool) {
	dir, file := filepath.Split(path)
	vbs := dir + strings.Replace(file, VboxExtension, SummaryExtension, 1)
	fmt.Printf("Looking for summary in %v\n", vbs)

	if _, err := os.Stat(vbs); err != nil && os.IsNotExist(err) {
		return vbs, false
	}
	return vbs, true
}

func loadSummary(path string) (*FileSummary, error) {
	var summary FileSummary
	file, err := os.Open(path)
	if err != nil {
		return &summary, err
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		return &summary, err
	}

	json.Unmarshal(bytes, &summary)
	fmt.Printf("Loaded summary from %v\n", path)
	return &summary, nil
}

func saveSummary(path string, summary *FileSummary) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("Failed to create summary at %v - %v\n", f, err)
	}
	defer f.Close()

	data, err := json.Marshal(summary)

	if err != nil {
		fmt.Printf("Failed to marshal summary to json at %v - %v\n", f, err)
	}

	if _, err := f.Write(data); err != nil {
		fmt.Printf("Failed to write summary at %v - %v\n", f, err)
	}
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
