package vboxanalyser

import (
	"errors"
	"fmt"
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
	Rows []VboFileDataRow
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

func (file *VboFile) AppendDataRow(row VboFileDataRow) {
	file.Data.Rows = append(file.Data.Rows, row)
}

func NewVboFileDataRow(data []interface{}) VboFileDataRow {
	return VboFileDataRow{data}
}

func (file *VboFile) MaxValueWithFunc(extractor func(*VboFileDataRow) (interface{}, error)) (float64, error) {
	var max float64

	for r := range file.Data.Rows {
		val, err := extractor(&file.Data.Rows[r])
		if err != nil {
			return 0, err
		}
		if val.(float64) > max {
			max = val.(float64)
		}
	}
	return max, nil
}

func ExtractValueFunctionFactory(channel string, file *VboFile) func(*VboFileDataRow) (interface{}, error) {
	return func(r *VboFileDataRow) (interface{}, error) {
		if v, ok := file.Columns[channel]; ok {
			return r.data[v], nil
		}
		return nil, errors.New(fmt.Sprintf("No channel name %v", channel))
	}
}
