package main

import (
	"math"
)

type Point struct {
	Val float64
	Ts  uint32
}

var nanPoint = Point{Val: math.NaN(), Ts: 0}

func cleanTimestamp(timestamp uint32, interval uint32) uint32 {
	rem := timestamp % interval
	if rem == 0 {
		return timestamp
	}
	result := (timestamp/interval)*interval + interval
	return result
}

func Fix(in []Point, from, to, interval uint32) []Point {

	from = cleanTimestamp(from, interval)
	to = cleanTimestamp(to, interval)
	outSize := (to - from) / interval
	out := make([]Point, outSize)

	offIdx := 0
	found := false

	for timestamp, outIdx := from, 0; timestamp < to; timestamp, outIdx = timestamp+interval, outIdx+1 {
		for inpIdx := offIdx; inpIdx < len(in); inpIdx++ {
			in[inpIdx].Ts = cleanTimestamp(in[inpIdx].Ts, interval)
			if in[inpIdx].Ts == timestamp {
				found = true
				out[outIdx] = in[inpIdx]
				offIdx = inpIdx
				break
			}
		}

		if found {
			found = false
		} else {
			nanPoint.Ts = timestamp
			out[outIdx] = nanPoint
		}
	}

	return out
}
