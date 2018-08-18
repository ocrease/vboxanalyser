package vbo

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	timeLayout = `150405.00`
	dateLayout = "02/01/2006 @ 15:04:05"
)

var dateFormat = regexp.MustCompile(`(0[1-9]|[12]\d|3[01])/(0[1-9]|1[0-2])/([12]\d{3}) @ ([01]\d|2[0-3]):([0-5]\d):([0-5]\d)`)

//ParseFile creates a VboFile representation
func ParseFile(path string) File {
	data, _ := os.Open(path)
	defer data.Close()

	vboFile := File{Path: path, Columns: make(map[string]int)}

	section := "pre"

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
			vboFile.Data = &Data{MaxValues: make([]float64, len(vboFile.Columns))}
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
	vboFile.Laps = processLaps(vboFile)
	return vboFile
}

func processRow(section string, line string, vboFile *File) {
	switch section {
	case "pre":
		ds := dateFormat.FindString(line)
		if len(ds) > 0 {
			if date, err := time.Parse(dateLayout, ds); err == nil {
				vboFile.CreationTime = date
			}
		}
	case "data":
		vboFile.createDataRow(strings.Fields(line))
	// case "header":
	// 	vboFile.header = append(vboFile.header, strings.Fields(line)...)
	case "laptiming":
		fields := strings.Fields(line)
		if fields[0] == "Start" {
			lon1, _ := strconv.ParseFloat(fields[1], 64)
			lat1, _ := strconv.ParseFloat(fields[2], 64)
			vboFile.Start = LatLng{lat1 / 60, lon1 * -1 / 60}
		}
	case "columns":
		for i, v := range strings.Fields(line) {
			vboFile.Columns[v] = i
		}
	}
}

func (f *File) createDataRow(fields []string) {
	fieldIndex := f.Columns
	data := make([]interface{}, len(fieldIndex))

	for name, index := range fieldIndex {
		switch name {
		case "sats":
			data[index], _ = strconv.Atoi(fields[index])
		case "time":
			data[index] = fields[index]
		default:
			if v, err := strconv.ParseFloat(fields[index], 64); err == nil {
				data[index] = v
				f.Data.updateMaxValue(index, v)
			}
		}
	}

	f.Data.Rows = append(f.Data.Rows, DataRow{data})

}

func (data *Data) updateMaxValue(index int, val float64) {
	if cur := data.MaxValues[index]; val > cur {
		data.MaxValues[index] = val
	}
}

func (f *File) maxValue(channel string) (float64, error) {
	i, ok := f.Columns[channel]

	if !ok {
		return 0, fmt.Errorf("No channel name %v", channel)
	}

	return f.Data.MaxValues[i], nil
}

func (f *File) duration(s, e int) time.Duration {
	i, ok := f.Columns["time"]
	if !ok {
		return 0
	}
	rows := f.Data.Rows
	start := parseTime(rows[s].data[i].(string))
	end := parseTime(rows[e].data[i].(string))
	return end.Sub(start)
}

func parseTime(t string) time.Time {
	ts, err := time.Parse(timeLayout, t)
	if err != nil {
		fmt.Printf("Failed to parse duration %v - %v", t, err)
		return *new(time.Time)
	}
	return ts
}

func (f *File) numLaps() (completedLaps int) {
	for _, l := range f.Laps {
		if !l.Partial {
			completedLaps = completedLaps + 1
		}
	}
	return
}

func (f *File) fastestLap() jsonDuration {
	var d jsonDuration
	for _, l := range f.Laps {
		if !l.Partial {
			if d == 0 || d > l.LapTime {
				d = l.LapTime
			}
		}
	}
	return d
}
