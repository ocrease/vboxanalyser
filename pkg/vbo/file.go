package vbo

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

//ParseFile creates a VboFile representation
func ParseFile(path string) File {
	data, _ := os.Open(path)
	defer data.Close()

	vboFile := File{Path: path, Columns: make(map[string]int)}

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
	return vboFile
}

func processRow(section string, line string, vboFile *File) {
	switch section {
	case "data":
		vboFile.CreateDataRow(strings.Fields(line))
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