package models

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang/geo/s1"

	"github.com/golang/geo/s2"
)

const (
	EarthRadius float64 = 6367000.0
)

var gateWidth = s1.ChordAngleFromAngle(s1.Angle(12.5 / EarthRadius))

//VboFile Main file object
type VboFile struct {
	path         string
	creationTime string
	header       []string
	comments     VboFileComments
	start        s2.LatLng
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
		case "[laptiming]":
			section = "laptiming"
		default:
			if strings.HasPrefix(line, "[") {
				section = "unknown"
			} else if len(line) > 0 {
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
	case "laptiming":
		fields := strings.Fields(line)
		if fields[0] == "Start" {
			lon1, _ := strconv.ParseFloat(fields[1], 64)
			lat1, _ := strconv.ParseFloat(fields[2], 64)
			vboFile.start = s2.LatLngFromDegrees(lat1/60, lon1*-1/60)
		}
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

	case "temp":
		f = func(r *VboFileDataRow) float64 {
			return r.coolantTemp
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
	data                 []string
	sats                 int
	time                 time.Time
	latLon               s2.LatLng
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

	var lat, lon float64

	for name, index := range fieldIndex {
		switch name {
		case "sats":
			row.sats, _ = strconv.Atoi(fields[index])
		case "lat":
			lat, _ = strconv.ParseFloat(fields[index], 64)
		case "long":
			lon, _ = strconv.ParseFloat(fields[index], 64)
		case "Corretced_coolant_temp":
			row.coolantTemp, _ = strconv.ParseFloat(fields[index], 64)
		case "LOT_Engine_Spd":
			row.engineSpeed, _ = strconv.ParseFloat(fields[index], 64)
		case "velocity":
			row.velocity, _ = strconv.ParseFloat(fields[index], 64)
		case "Veh_Spd_MPH":
			row.velocityMph, _ = strconv.ParseFloat(fields[index], 64)
		}
	}

	row.latLon = s2.LatLngFromDegrees(lat/60, lon*-1/60)
	return row

}

func (file *VboFile) NumLaps() int {
	rows := file.data.rows
	var numLaps = 0
	var crossingStartLine = false
	for i, v := range rows {
		if i > 0 {
			start := s2.PointFromLatLng(rows[i-1].latLon)
			finish := s2.PointFromLatLng(v.latLon)
			inRange := s2.IsDistanceLess(s2.PointFromLatLng(file.start), start, finish, gateWidth)
			if !inRange {
				crossingStartLine = false
			} else if !crossingStartLine {
				crossingStartLine = true
				numLaps = numLaps + 1
			}
		}
	}

	return numLaps

}

func (file *VboFile) Distance() float64 {
	//latLons := make([]s2.LatLng, len(file.data.rows))
	var distance float64
	for i, v := range file.data.rows {
		if i == 0 {
			continue
		}
		distance = distance + v.latLon.Distance(file.data.rows[i-1].latLon).Radians()
		//latLons = append(latLons, v.latLon)
	}
	//polyline := s2.PolylineFromLatLngs(latLons)
	return distance * (EarthRadius / 1000)
}
