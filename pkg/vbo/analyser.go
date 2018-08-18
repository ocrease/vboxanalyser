package vbo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Masterminds/semver"
)

type Analyser struct{}

const (
	vboxExtension    = ".vbo"
	summaryExtension = ".vbs"
	summaryVersion   = "0.2.0"
)

var summaryVersionConstraint, _ = semver.NewConstraint(summaryVersion)

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

			if filepath.Ext(path) == vboxExtension {
				handleFile(path, info, consumer)
			}
		}
		return nil
	}
}

func handleFile(path string, info os.FileInfo, consumer func(FileSummary)) {
	vbs, vbsExists := summaryExists(path)
	if vbsExists {
		summary, err := loadSummary(vbs)
		var valid bool
		if err == nil {
			valid = validateSummary(info, summary)
			if valid {
				consumer(*summary)
				return
			}
		}
	}

	fmt.Printf("Analysing %v\n", path)
	file := ParseFile(path)
	summary := generateSummary(&file, info)
	saveSummary(vbs, summary)
	consumer(*summary)
}

func generateSummary(file *File, info os.FileInfo) *FileSummary {
	rpm, err := file.MaxValue("LOT_Engine_Spd")
	if err != nil {
		rpm = 0
	}
	vel, err := file.MaxValue("velocity")
	if err != nil {
		vel = 0
	}
	return &FileSummary{
		Version:      summaryVersion,
		CreationTime: file.CreationTime,
		ModTime:      info.ModTime(),
		Path:         file.Path,
		NumLaps:      NumLaps(file),
		MaxVelocity:  vel,
		MaxRpm:       rpm}
}

func summaryExists(path string) (string, bool) {
	dir, file := filepath.Split(path)
	vbs := dir + strings.Replace(file, vboxExtension, summaryExtension, 1)

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

	return &summary, nil
}

func validateSummary(info os.FileInfo, summary *FileSummary) bool {
	version, err := semver.NewVersion(summary.Version)
	if err != nil {
		return false
	}
	if !summaryVersionConstraint.Check(version) {
		return false
	}
	if info.ModTime().After(summary.ModTime) {
		return false
	}
	return true
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

type FileSummary struct {
	Version      string    `json:"version,omitempty"`
	CreationTime time.Time `json:"creationtime"`
	ModTime      time.Time `json:"modtime"`
	Path         string    `json:"path"`
	NumLaps      uint32    `json:"numlaps"`
	MaxVelocity  float64   `json:"maxvelocity"`
	MaxRpm       float64   `json:"maxrpm"`
}
