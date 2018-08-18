package vbo

import (
	"time"
)

//File Main file object
type File struct {
	Path         string
	CreationTime time.Time
	comments     Comments
	Start        LatLng
	Columns      map[string]int
	Data         *Data
	Laps         []Lap
}

// type VboFileChannelUnits struct {
// }

//Comments contains file information
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

//Data contains all the data rows
type Data struct {
	Rows      []DataRow
	MaxValues []float64
}

type Lap struct {
	startIndex int
	endIndex   int
	Partial    bool         `json:"partial"`
	LapTime    jsonDuration `json:"laptime"`
	maxValues  []float64
}

//DataRow contains the data fields in a row
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
