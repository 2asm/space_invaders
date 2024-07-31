//go:build js && wasm

package game

import (
	"math/rand"
	"time"
)

type alienArmy struct {
	aliens     []*alien
	aliveCount int
	direction  direction
}

const (
	boxHeight       = 10
	boxWidth        = 14
	startX          = 15
	startY          = 9
	armyRowCount    = 5
	armyColumnCount = 11
)

func initialArmy() *alienArmy {
	ret := &alienArmy{aliens: make([]*alien, 0)}
	kind := _I3
	id := 0
	for row := range armyRowCount {
		if row == 1 {
			kind = _I2
		} else if row == 3 {
			kind = _I1
		}
		for col := range armyColumnCount {
			x := startX + boxHeight*row + alienXOffset[kind]
			y := startY + boxWidth*col + alienYOffset[kind]
			ret.aliens = append(ret.aliens, newAlien(id, kind, x, y))
			id += 1
		}
	}
	ret.direction = _RIGHT
	ret.aliveCount = armyRowCount * armyColumnCount
	return ret
}

// aliend missile launch delay int
// (milliseconds)
func (aa *alienArmy) getLaunchDelay() time.Duration {
	r := rand.Intn(aa.aliveCount+1) + aa.aliveCount
	return time.Duration(r) * time.Millisecond * 100
}

func (aa *alienArmy) moveLeft(steps int) {
	for _, a := range aa.aliens {
		if !a.died {
			a.switchState()
			a.moveLeft(steps)
		}
	}
}

func (aa *alienArmy) moveRight(steps int) {
	for _, a := range aa.aliens {
		if !a.died {
			a.switchState()
			a.moveRight(steps)
		}
	}
}

func (aa *alienArmy) move(steps int) {
	// if aa.rightMostY() >= canvasDimY-1 || aa.leftMostY() <= 0 {
	if aa.rightMostY() >= canvasDimY-1-5 || aa.leftMostY() <= 0+5 {
		aa.changeDirection()
	}
	if aa.direction == _RIGHT {
		aa.moveRight(steps)
	} else {
		aa.moveLeft(steps)
	}
}

func (aa *alienArmy) render() {
	for _, a := range aa.aliens {
		if !a.died {
			a.render(a.state)
		}
	}
}

func (aa *alienArmy) clear() {
	for _, a := range aa.aliens {
		if !a.died {
			a.clear(a.state)
		}
	}
}

func (aa *alienArmy) changeDirection() {
	aa.direction = aa.direction.opposite()
}

func (aa *alienArmy) leftMostY() int {
	Y := -1
	for idx := range aa.aliens {
		if !aa.aliens[idx].died {
			if Y == -1 {
				Y = aa.aliens[idx].baseY
			} else {
				Y = min(Y, aa.aliens[idx].baseY)
			}
		}
	}
	return Y
}

func (aa *alienArmy) rightMostY() int {
	Y := -1
	for _, a := range aa.aliens {
		if !a.died {
			Y = max(Y, a.baseY+alienWidth[a.kind]-1)
		}
	}
	return Y
}
