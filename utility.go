package main

import "math"

func Distance(p1, p2 Point) int {
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y
	dz := p2.Z - p1.Z
	return int(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}
