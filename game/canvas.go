//go:build js && wasm

package game

import "syscall/js"

const ( // dim
	canvasDimX = 150
	canvasDimY = 250
)

var (
	scale         = 4
	gameCanvas    js.Value
	missileCanvas js.Value
)

func init() {
	// todo: change scale accoding to resolution
	c := js.Global().Get("document").Call("getElementById", "gameCanvas")
	c.Set("height", scale*canvasDimX)
	c.Set("width", scale*canvasDimY)
	gameCanvas = c.Call("getContext", "2d")

	c2 := js.Global().Get("document").Call("getElementById", "missileCanvas")
	c2.Set("height", scale*canvasDimX)
	c2.Set("width", scale*canvasDimY)
	missileCanvas = c2.Call("getContext", "2d")

	c3 := js.Global().Get("document").Call("getElementById", "backCanvas")
	c3.Set("height", scale*canvasDimX)
	c3.Set("width", scale*canvasDimY)
}
