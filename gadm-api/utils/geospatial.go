package utils

type Point struct {
	Lng float64
	Lat float64
}

func NewPointLngLat(lng float64, lat float64) Point {
	return Point{
		Lng: lng,
		Lat: lat,
	}
}
