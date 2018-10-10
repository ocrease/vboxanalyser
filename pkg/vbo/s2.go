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
	startLine := s2.PointFromLatLng(s2.LatLngFromDegrees(f.Start.Lat, f.Start.Lng))
	rows := f.Data.Rows

	latIndex := f.Columns[Latitude]
	longIndex := f.Columns[Longitude]
	latLng := func(r DataRow) s2.LatLng {
		s1Lat := r.GetValue(latIndex).(float64)
		s1Lng := r.GetValue(longIndex).(float64)
		return s2.LatLngFromDegrees(s1Lat/60, s1Lng*-1/60)
	}
	startFraction := func(i int) time.Duration {
		if i <= 1 || i == len(rows)-1 {
			return 0
		}
		a := s2.PointFromLatLng(latLng(rows[i-1]))
		b := s2.PointFromLatLng(latLng(rows[i]))
		x := s2.Project(startLine, a, b)
		sf := s2.DistanceFraction(x, a, b)
		sf = 1.0 - sf

		duration := float64(f.duration(i-1, i)) * sf
		return time.Duration(duration)
	}

	duration := func(s int, e int) time.Duration {
		sd := startFraction(s)
		ed := startFraction(e)
		return f.duration(s, e) + sd - ed
	}

	start := 0
	end := 0
	var s s2.Point
	var e s2.Point
	minDistance := gateWidth
	partial := true
	firstLap := true
	inRange := false
	startedInRange := false
	nearStartLine := false
	for i, v := range rows {
		if i == 0 {
			e = s2.PointFromLatLng(latLng(v))
			comp := s2.CompareDistance(
				startLine,
				e,
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
		e = s2.PointFromLatLng(latLng(v))
		inRange = s2.IsDistanceLess(
			startLine,
			s,
			e,
			gateWidth)
		if inRange {
			d, updated := s2.UpdateMinDistance(
				startLine,
				s,
				e,
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
			laps = append(laps, newLap(f, start, end, partial, duration))
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
		laps = append(laps, newLap(f, start, end, partial, duration))
		start = end
	}
	if start != (len(rows) - 1) {
		laps = append(laps, newLap(f, start, len(rows)-1, true, duration))
	}
	return
}

func newLap(f File, s int, e int, partial bool, duration func(a, b int) time.Duration) Lap {
	return Lap{startIndex: s, endIndex: e, Partial: partial, LapTime: jsonDuration(duration(s, e).Round(10 * time.Millisecond))}
}

func Distance(f *File) float64 {
	//latLons := make([]s2.LatLng, len(file.data.rows))
	latIndex := f.Columns[Latitude]
	longIndex := f.Columns[Longitude]
	latLng := func(r DataRow) s2.LatLng {
		s1Lat := r.GetValue(latIndex).(float64)
		s1Lng := r.GetValue(longIndex).(float64)
		return s2.LatLngFromDegrees(s1Lat/60, s1Lng*-1/60)
	}
	var distance float64
	for i, v := range f.Data.Rows {
		if i == 0 {
			continue
		}
		distance = distance + latLng(v).Distance(latLng(f.Data.Rows[i-1])).Radians()
		//latLons = append(latLons, v.latLon)
	}
	//polyline := s2.PolylineFromLatLngs(latLons)
	return distance * (earthRadius / 1000)
}
