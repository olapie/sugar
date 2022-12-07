package types

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/geo/s2"
)

const (
	PI          = 3.141_592_65
	EarthRadius = 6_378_100
	EarthCircle = 2 * PI * EarthRadius
	Degree      = EarthCircle * 1000 / 360
)

func NewPoint(x, y float64) *Point {
	return &Point{X: x, Y: y}
}

func NewPointFromString(s string) (*Point, error) {
	l := strings.Split(s, ",")
	if len(l) != 2 {
		return nil, errors.New("invalid format")
	}
	x := strings.TrimSpace(l[0])
	if x == "" {
		return nil, errors.New("x is empty")
	}
	p := new(Point)
	var err error
	p.X, err = strconv.ParseFloat(x, 64)
	if err != nil {
		return nil, fmt.Errorf("parse x %s: %w", x, err)
	}
	y := strings.TrimSpace(l[1])
	if y == "" {
		return nil, errors.New("y is empty")
	}
	p.Y, err = strconv.ParseFloat(y, 64)
	if err != nil {
		return nil, fmt.Errorf("parse y  %s: %w", x, err)
	}
	return p, nil
}

func (p *Point) Distance(v *Point) int {
	p1 := s2.PointFromLatLng(s2.LatLngFromDegrees(p.Y, p.X))
	p2 := s2.PointFromLatLng(s2.LatLngFromDegrees(v.Y, v.X))
	d := p1.Distance(p2)
	return int(d.Radians() * EarthRadius)
}
