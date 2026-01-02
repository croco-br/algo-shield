package rules

import (
	"encoding/json"
)

// Point represents a geographic coordinate with latitude and longitude.
type Point struct {
	Lat float64
	Lon float64
}

// ParsePolygon parses a JSON string of polygon coordinates.
// Expected format: [[lat1, lon1], [lat2, lon2], ...]
// Returns a slice of Points representing the polygon vertices.
func ParsePolygon(jsonStr string) ([]Point, error) {
	var coords [][]float64
	if err := json.Unmarshal([]byte(jsonStr), &coords); err != nil {
		return nil, err
	}

	polygon := make([]Point, 0, len(coords))
	for _, coord := range coords {
		if len(coord) != 2 {
			continue
		}
		polygon = append(polygon, Point{Lat: coord[0], Lon: coord[1]})
	}

	return polygon, nil
}

// GetCoordinatesFromMetadata extracts lat/lon from transaction metadata.
// Supports multiple formats:
//   - metadata.latitude, metadata.longitude (float64)
//   - metadata.lat, metadata.lon (float64)
//   - metadata.location.lat, metadata.location.lon (nested)
//   - metadata.location.latitude, metadata.location.longitude (nested)
func GetCoordinatesFromMetadata(metadata map[string]any) (lat, lon float64, ok bool) {
	if metadata == nil {
		return 0, 0, false
	}

	// Try latitude/longitude
	if latVal, exists := metadata["latitude"]; exists {
		if lonVal, exists := metadata["longitude"]; exists {
			lat, latOk := toFloat64(latVal)
			lon, lonOk := toFloat64(lonVal)
			if latOk && lonOk {
				return lat, lon, true
			}
		}
	}

	// Try lat/lon
	if latVal, exists := metadata["lat"]; exists {
		if lonVal, exists := metadata["lon"]; exists {
			lat, latOk := toFloat64(latVal)
			lon, lonOk := toFloat64(lonVal)
			if latOk && lonOk {
				return lat, lon, true
			}
		}
	}

	// Try nested location object
	if location, exists := metadata["location"]; exists {
		if locMap, ok := location.(map[string]any); ok {
			// Try lat/lon in location
			if latVal, exists := locMap["lat"]; exists {
				if lonVal, exists := locMap["lon"]; exists {
					lat, latOk := toFloat64(latVal)
					lon, lonOk := toFloat64(lonVal)
					if latOk && lonOk {
						return lat, lon, true
					}
				}
			}
			// Try latitude/longitude in location
			if latVal, exists := locMap["latitude"]; exists {
				if lonVal, exists := locMap["longitude"]; exists {
					lat, latOk := toFloat64(latVal)
					lon, lonOk := toFloat64(lonVal)
					if latOk && lonOk {
						return lat, lon, true
					}
				}
			}
		}
	}

	return 0, 0, false
}

// PointInPolygon checks if a point is inside a polygon using the ray casting algorithm.
// This is a standard algorithm that counts how many times a ray from the point
// crosses the polygon boundary. If the count is odd, the point is inside.
func PointInPolygon(lat, lon float64, polygon []Point) bool {
	n := len(polygon)
	if n < 3 {
		return false
	}

	inside := false
	j := n - 1

	for i := 0; i < n; i++ {
		// Check if the ray from (lat, lon) going right crosses the edge from polygon[i] to polygon[j]
		if ((polygon[i].Lon > lon) != (polygon[j].Lon > lon)) &&
			(lat < (polygon[j].Lat-polygon[i].Lat)*(lon-polygon[i].Lon)/(polygon[j].Lon-polygon[i].Lon)+polygon[i].Lat) {
			inside = !inside
		}
		j = i
	}

	return inside
}

// toFloat64 converts various numeric types to float64.
// Returns the converted value and true if successful, or 0 and false if not.
func toFloat64(v any) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case int32:
		return float64(val), true
	default:
		return 0, false
	}
}
