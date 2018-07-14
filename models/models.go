package models

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

//VboFile Main file object
type VboFile struct {
	path         string
	creationTime string
	header       []string
	comments     VboFileComments
	columns      map[string]int
	data         *VboFileData
}

//ParseFile creates a VboFile representation
func ParseFile(path string) VboFile {
	data, _ := os.Open(path)
	defer data.Close()

	vboFile := VboFile{path: path, columns: make(map[string]int), data: &VboFileData{}}

	var section string

	scanner := bufio.NewScanner(data)

	for scanner.Scan() {
		line := scanner.Text()
		switch line {
		case "[header]":
			section = "header"
		case "[comments]":
			section = "comments"
		case "[column names]":
			section = "columns"
		case "[data]":
			section = "data"
		default:
			if strings.HasPrefix(line, "[") {
				section = "unknown"
			} else {
				processRow(section, line, &vboFile)
			}
		}

	}
	return vboFile
}

func processRow(section string, line string, vboFile *VboFile) {
	switch section {
	case "data":
		row := NewVboFileDataRow(strings.Fields(line), vboFile.columns)
		vboFile.appendDataRow(row)
	case "header":
		vboFile.header = append(vboFile.header, strings.Fields(line)...)
	case "columns":
		for i, v := range strings.Fields(line) {
			vboFile.columns[v] = i
		}
	}
}

func (file *VboFile) appendDataRow(row VboFileDataRow) {
	file.data.rows = append(file.data.rows, row)
}

//NumDataPoints returns the number of data points
func (file *VboFile) NumDataPoints() int {
	return len(file.data.rows)
}

//NumColumns returns the number of columns in the file
func (file *VboFile) NumColumns() int {
	return len(file.columns)
}

func (file *VboFile) MaxValueWithFunc(extractor func(*VboFileDataRow) float64) float64 {
	var max float64

	for r := range file.data.rows {
		val := extractor(&file.data.rows[r])
		if val > max {
			max = val
		}
	}
	return max
}

func ExtractValueFunctionFactory(channel string) func(*VboFileDataRow) float64 {
	var f func(*VboFileDataRow) float64

	switch channel {
	case "rpm":
		f = func(r *VboFileDataRow) float64 {
			return r.engineSpeed
		}

	case "speedKph":
		f = func(r *VboFileDataRow) float64 {
			return r.velocity
		}

	case "speedMph":
		f = func(r *VboFileDataRow) float64 {
			return r.velocityMph
		}
	}
	return f
}

// type VboFileChannelUnits struct {
// }

//VboFileComments contains file information
type VboFileComments struct {
	vboxVersion  string
	serialNumber string
	engine       string
	gps          string
	rate         float32
	maxRate      float32
}

// type VboFileLapTiming struct {
// }

// type VboFileAvi struct {
// }

//VboFileData contains all the data rows
type VboFileData struct {
	rows []VboFileDataRow
}

//VboFileDataRow contains the data fields in a row
type VboFileDataRow struct {
	sats                 int
	time                 time.Time
	lat                  string
	lon                  string
	velocity             float64
	heading              float64
	height               float64
	vertVelocity         float64
	engineSpeed          float64
	coolantTemp          float64
	velocityMph          float64
	correctedCoolantTemp float64
}

func NewVboFileDataRow(fields []string, fieldIndex map[string]int) VboFileDataRow {
	row := VboFileDataRow{}

	for name, index := range fieldIndex {
		switch name {
		case "sats":
			row.sats, _ = strconv.Atoi(fields[index])
		case "LOT_Engine_Spd":
			row.engineSpeed, _ = strconv.ParseFloat(fields[index], 64)
		case "velocity":
			row.velocity, _ = strconv.ParseFloat(fields[index], 64)
		case "Veh_Spd_MPH":
			row.velocityMph, _ = strconv.ParseFloat(fields[index], 64)
		}
	}
	return row

}
