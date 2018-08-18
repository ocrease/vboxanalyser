package vbo

import (
	"time"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

const (
	earthRadius float64 = 6367000.0
)

var gateWidth = s1.ChordAngleFromAngle(s1.Angle(12.5 / earthRadius))

func processLaps(f File) (laps []Lap) {
	startLine := s2.LatLngFromDegrees(f.Start.Lat, f.Start.Lng)
	rows := f.Data.Rows

	latIndex := f.Columns["lat"]
	longIndex := f.Columns["long"]
	latLng := func(r DataRow) s2.LatLng {
		s1Lat := r.GetValue(latIndex).(float64)
		s1Lng := r.GetValue(longIndex).(float64)
		return s2.LatLngFromDegrees(s1Lat/60, s1Lng*-1/60)
	}

	start := 0
	end := 0
	var s s2.LatLng
	var e s2.LatLng
	minDistance := gateWidth
	partial := true
	firstLap := true
	inRange := false
	startedInRange := false
	nearStartLine := false
	for i, v := range rows {
		if i == 0 {
			e = latLng(v)
			comp := s2.CompareDistance(
				s2.PointFromLatLng(startLine),
				s2.PointFromLatLng(e),
				gateWidth)
			if comp < 1 {
				inRange = true
				nearStartLine = true
				startedInRange = true
				minDistance = s1.ChordAngleFromAngle(startLine.Distance(e))
			}
			continue
		}
		s = e
		e = latLng(v)
		inRange = s2.IsDistanceLess(
			s2.PointFromLatLng(startLine),
			s2.PointFromLatLng(s),
			s2.PointFromLatLng(e),
			gateWidth)
		if inRange {
			d, updated := s2.UpdateMinDistance(
				s2.PointFromLatLng(startLine),
				s2.PointFromLatLng(s),
				s2.PointFromLatLng(e),
				minDistance)
			if updated {
				end = i
				minDistance = d
				startedInRange = false
			}
			nearStartLine = true
		} else if nearStartLine { // previously near the start line
			if firstLap {
				partial = !startedInRange
			}
			laps = append(laps, newLap(start, end, partial, f.duration(start, end)))
			start = end
			partial = false
			firstLap = false
			minDistance = gateWidth
			nearStartLine = false
		}
	}
	if nearStartLine && start != end {
		if firstLap {
			partial = !startedInRange
		}
		laps = append(laps, newLap(start, end, partial, f.duration(start, end)))
		start = end
	}
	if start != (len(rows) - 1) {
		laps = append(laps, newLap(start, len(rows)-1, true, f.duration(start, len(rows)-1)))
	}
	return
}

func newLap(s int, e int, partial bool, lapTime time.Duration) Lap {
	return Lap{startIndex: s, endIndex: e, Partial: partial, LapTime: jsonDuration(lapTime)}
}

// func Distance(f *vbo.File) float64 {
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
