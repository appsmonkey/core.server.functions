package zones

import "math"

// Point point
type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Zone representing a complete zone with all needed points
type Zone struct {
	Title    string   `json:"title"`
	Polygons []*Point `json:"polygon"`
}

var zones map[string]Zone

func init() {
	zones = make(map[string]Zone, 0)
}

// AddZones will add zones to the internal storage
func AddZones(in []Zone) {
	if len(in) == 0 {
		return
	}

	for _, z := range in {
		zones[z.Title] = z
	}
}

func (z *Zone) IsClosed() bool {
	if len(z.Polygons) < 3 {
		return false
	}

	return true
}

// Contains returns whether or not the current Polygon contains the passed in Point.
func (z *Zone) Contains(point *Point) bool {
	if !z.IsClosed() {
		return false
	}

	start := len(z.Polygons) - 1
	end := 0

	contains := intersectsWithRaycast(point, z.Polygons[start], z.Polygons[end])

	for i := 1; i < len(z.Polygons); i++ {
		if intersectsWithRaycast(point, z.Polygons[i-1], z.Polygons[i]) {
			contains = !contains
		}
	}

	return contains
}

// Using the raycast algorithm, this returns whether or not the passed in point
// Intersects with the edge drawn by the passed in start and end points.
// Original implementation: http://rosettacode.org/wiki/Ray-casting_algorithm#Go
func intersectsWithRaycast(point *Point, start *Point, end *Point) bool {
	// Always ensure that the the first point
	// has a y coordinate that is less than the second point
	if start.Lng > end.Lng {
		// Switch the points if otherwise.
		start, end = end, start
	}

	// Move the point's y coordinate
	// outside of the bounds of the testing region
	// so we can start drawing a ray
	for point.Lng == start.Lng || point.Lng == end.Lng {
		newLng := math.Nextafter(point.Lng, math.Inf(1))
		point = &Point{Lat: point.Lat, Lng: newLng}
	}

	// If we are outside of the polygon, indicate so.
	if point.Lng < start.Lng || point.Lng > end.Lng {
		return false
	}

	if start.Lat > end.Lat {
		if point.Lat > start.Lat {
			return false
		}
		if point.Lat < end.Lat {
			return true
		}

	} else {
		if point.Lat > end.Lat {
			return false
		}
		if point.Lat < start.Lat {
			return true
		}
	}

	raySlope := (point.Lng - start.Lng) / (point.Lat - start.Lat)
	diagSlope := (end.Lng - start.Lng) / (end.Lat - start.Lat)

	return raySlope >= diagSlope
}

// ZoneByPoint will return the zone that encapsulates the provided point
func ZoneByPoint(p *Point) *Zone {
	// No point provided or no zones loaded
	if p == nil || len(zones) == 0 {
		return nil
	}

	for _, z := range zones {
		if z.Contains(p) {
			return &z
		}
	}

	// If we did not find any zone that encapsulates the point then just return nil
	return nil
}

// ZoneByName will return the zone represented by that name or `nil`
func ZoneByName(n string) *Zone {
	// No point provided or no zones loaded
	if len(n) == 0 || len(zones) == 0 {
		return nil
	}

	z, ok := zones[n]
	if !ok {
		return nil
	}

	// If we did not find any zone by that name then just return nil
	return &z
}

// pointInZone will calculate if the provided point is encapsulated by the provided zone
// func pointInZone(p *Point, z *Zone) bool {
// 	polyCorners := len(z.Polygons)
// 	polyX := make([]float64, 0)
// 	polyY := make([]float64, 0)
// 	var x = p.Lat
// 	var y = p.Lng
// 	var j = polyCorners - 1

// 	for _, pol := range z.Polygons {
// 		polyX = append(polyX, pol.Lat)
// 		polyY = append(polyY, pol.Lng)
// 	}

// 	for i := 0; i < polyCorners; i++ {
// 		if (polyY[i] < y && polyY[j] >= y || polyY[j] < y && polyY[i] >= y) && (polyX[i] <= x || polyX[j] <= x) {
// 			if polyX[i]+(y-polyY[i])/(polyY[j]-polyY[i])*(polyX[j]-polyX[i]) < x {
// 				return true
// 			}
// 		}
// 	}

// 	return false
// }
