//go:build js && wasm

package game

type coord struct {
	x, y int
}

func newCoord(x, y int) coord {
	return coord{x, y}
}

func (c coord) add(x, y int) coord {
	return newCoord(c.x+x, c.y+y)
}
