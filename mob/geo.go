package mob

import (
	"code.olapie.com/sugar/v2/types"
)

type Point types.Point

func NewPoint() *Point {
	return new(Point)
}

type Place types.Place

func NewPlace() *Place {
	return new(Place)
}

func (p *Place) SetCoordinate(c *Point) {
	p.Coordinate = (*types.Point)(c)
}

func (p *Place) GetCoordinate() *Point {
	return (*Point)(p.Coordinate)
}
