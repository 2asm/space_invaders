//go:build js && wasm

package game

type player struct {
	baseX, baseY int
}

var (
	playerColor = "grey"
	basePlayer  = []coord{
                                {0, 3},
                        {1, 2}, {1, 3}, {1, 4},
                {2, 1}, {2, 2}, {2, 3}, {2, 4}, {2, 5},
                {3, 1}, {3, 2}, {3, 3}, {3, 4}, {3, 5},
                {4, 1}, {4, 2}, {4, 3}, {4, 4}, {4, 5},
		{5, 0}, {5, 1}, {5, 2}, {5, 3}, {5, 4}, {5, 5}, {5, 6},
		{6, 0}, {6, 1}, {6, 2}, {6, 3}, {6, 4}, {6, 5}, {6, 6},
		{7, 0}, {7, 1},                         {7, 5}, {7, 6},
	}
)

func newPlayer(baseX, baseY int) *player {
	return &player{
		baseX: baseX,
		baseY: baseY,
	}
}

func initialPlayer() *player {
	return newPlayer(canvasDimX-15, canvasDimY/2)
}

func (p *player) getBaseCoordinates() []coord {
	return basePlayer
}

func (p *player) move(d direction) {
	p.clear()
	switch d {
	case _LEFT:
		p.baseY -= 1
	case _RIGHT:
		p.baseY += 1
	case _UP:
		p.baseX -= 1
	case _DOWN:
		p.baseX += 1
	default:
		panic("unreachable")
	}
	p.render()
}

func (p *player) launchMissile() *missile {
	return newMissile(p.baseX-1-3, p.baseY+2, playerMissile)
}

func (p *player) render() {
	gameCanvas.Set("fillStyle", playerColor)
	for _, c := range basePlayer {
		gameCanvas.Call("fillRect", (p.baseY+c.y)*scale, (p.baseX+c.x)*scale, scale, scale)
	}
	gameCanvas.Call("fill")
}

func (p *player) clear() {
	for _, c := range basePlayer {
		gameCanvas.Call("clearRect", (p.baseY+c.y)*scale, (p.baseX+c.x)*scale, scale, scale)
	}
	gameCanvas.Call("fill")
}
