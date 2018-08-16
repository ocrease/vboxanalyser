package vbo

import (
	"fmt"
	"strconv"
)

//VboFile Main file object
type File struct {
	Path         string
	CreationTime string
	comments     Comments
	Start        LatLng
	Columns      map[string]int
	Data         *Data
}

// type VboFileChannelUnits struct {
// }

//VboFileComments contains file information
type Comments struct {
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
type Data struct {
	Rows      []DataRow
	MaxValues []float64
}

//VboFileDataRow contains the data fields in a row
type DataRow struct {
	data []interface{}
}

func (r *DataRow) GetValue(index int) interface{} {
	return r.data[index]
}

type LatLng struct {
	Lat float64
	Lng float64
}

func (file *File) CreateDataRow(fields []string) {
	fieldIndex := file.Columns
	data := make([]interface{}, len(fieldIndex))

	for name, index := range fieldIndex {
		switch name {
		case "sats":
			data[index], _ = strconv.Atoi(fields[index])
		default:
			if v, err := strconv.ParseFloat(fields[index], 64); err == nil {
				data[index] = v
				file.Data.updateMaxValue(index, v)
			}
		}
	}

	file.Data.Rows = append(file.Data.Rows, DataRow{data})

}

func (data *Data) updateMaxValue(index int, val float64) {
	if cur := data.MaxValues[index]; val > cur {
		data.MaxValues[index] = val
	}
}

func (file *File) MaxValue(channel string) (float64, error) {
	i, ok := file.Columns[channel]

	if !ok {
		return 0, fmt.Errorf("No channel name %v", channel)
	}

	return file.Data.MaxValues[i], nil
}
