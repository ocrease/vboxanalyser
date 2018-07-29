package s2

import (
	"github.com/golang/geo/s1"
	"github.com/ocrease/vboxanalyser"

	"github.com/golang/geo/s2"
)

const (
	EarthRadius float64 = 6367000.0
)

var gateWidth = s1.ChordAngleFromAngle(s1.Angle(12.5 / EarthRadius))

func NumLaps(f *vboxanalyser.VboFile) int {
	rows := f.Data.Rows
	var numLaps = 0
	var crossingStartLine = false
	for i, v := range rows {
		if i > 0 {
			s1Lat := rows[i-1].GetValue(f.Columns["lat"]).(float64)
			s1Lng := rows[i-1].GetValue(f.Columns["long"]).(float64)
			s2Lat := v.GetValue(f.Columns["lat"]).(float64)
			s2Lng := v.GetValue(f.Columns["long"]).(float64)
			startLine := s2.LatLngFromDegrees(f.Start.Lat, f.Start.Lng)
			s := s2.LatLngFromDegrees(s1Lat/60, s1Lng*-1/60)
			f := s2.LatLngFromDegrees(s2Lat/60, s2Lng*-1/60)
			inRange := s2.IsDistanceLess(s2.PointFromLatLng(startLine), s2.PointFromLatLng(s), s2.PointFromLatLng(f), gateWidth)
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

// func Distance(f *vboxanalyser.VboFile) float64 {
// 	//latLons := make([]s2.LatLng, len(file.data.rows))
// 	var distance float64
// 	for i, v := range f.Data.Rows {
// 		if i == 0 {
// 			continue
// 		}
// 		distance = distance + v.latLon.Distance(file.data.rows[i-1].latLon).Radians()
// 		//latLons = append(latLons, v.latLon)
// 	}
// 	//polyline := s2.PolylineFromLatLngs(latLons)
// 	return distance * (EarthRadius / 1000)
// }
