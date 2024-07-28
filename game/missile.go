//go:build js && wasm

package game

type milssileKind int

const (
	alienMissile milssileKind = iota
	playerMissile
)

type missile struct {
	kind         milssileKind
	baseX, baseY int
	destroyed    bool
}

func newMissile(baseX, baseY int, kind milssileKind) *missile {
	return &missile{
		baseX: baseX,
		baseY: baseY,
		kind:  kind,
	}
}

var (
	baseMissile = []coord{
		{0, 0},
		{1, 0},
		{2, 0},
	}

	missileColor = []string{
		alienMissile:  "yellow",
		playerMissile: "red",
	}
)

func (m *missile) move() {
	switch m.kind {
	case alienMissile:
		m.baseX += 1
	case playerMissile:
		m.baseX -= 1
	default:
		panic("unreachable")
	}
}

func (m *missile) render() {
	missileCanvas.Set("fillStyle", missileColor[m.kind])
	for _, c := range baseMissile {
		missileCanvas.Call("fillRect", (m.baseY+c.y)*scale, (m.baseX+c.x)*scale, scale, scale)
	}
	missileCanvas.Call("fill")
}

func (m *missile) clear() {
	for _, c := range baseMissile {
		missileCanvas.Call("clearRect", (m.baseY+c.y)*scale, (m.baseX+c.x)*scale, scale, scale)
	}
	missileCanvas.Call("fill")
}