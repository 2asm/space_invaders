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
	baseMissile = [][]coord{
		playerMissile: {
			{0, 0},
			{1, 0},
			{2, 0},
		},
		alienMissile: {
			{0, 0},
			{1, 0},
			{2, 0},
		},
	}

	missileColor = []string{
		alienMissile:  "red",
		playerMissile: "yellow",
	}
)

func (m *missile) getBaseCoordinates() []coord {
	return baseMissile[m.kind]
}

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
	for _, c := range baseMissile[m.kind] {
		missileCanvas.Call("fillRect", (m.baseY+c.y)*scale, (m.baseX+c.x)*scale, scale, scale)
	}
	missileCanvas.Call("fill")
}

func (m *missile) clear() {
	for _, c := range baseMissile[m.kind] {
		missileCanvas.Call("clearRect", (m.baseY+c.y)*scale, (m.baseX+c.x)*scale, scale, scale)
	}
	missileCanvas.Call("fill")
}
