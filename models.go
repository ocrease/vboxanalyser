package vboxanalyser

import (
	"fmt"
	"strconv"
)

//VboFile Main file object
type VboFile struct {
	Path         string
	CreationTime string
	//header       []string
	comments VboFileComments
	Start    LatLng
	Columns  map[string]int
	Data     *VboFileData
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
	Rows      []VboFileDataRow
	MaxValues []float64
}

//VboFileDataRow contains the data fields in a row
type VboFileDataRow struct {
	data []interface{}
}

func (r *VboFileDataRow) GetValue(index int) interface{} {
	return r.data[index]
}

type LatLng struct {
	Lat float64
	Lng float64
}

func (file *VboFile) CreateDataRow(fields []string) {
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

	file.Data.Rows = append(file.Data.Rows, VboFileDataRow{data})

}

func (data *VboFileData) updateMaxValue(index int, val float64) {
	if cur := data.MaxValues[index]; val > cur {
		data.MaxValues[index] = val
	}
}

func (file *VboFile) MaxValue(channel string) (float64, error) {
	i, ok := file.Columns[channel]

	if !ok {
		return 0, fmt.Errorf("No channel name %v", channel)
	}

	return file.Data.MaxValues[i], nil
}
