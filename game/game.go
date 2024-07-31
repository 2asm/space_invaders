//go:build js && wasm

package game

import (
	"fmt"
	"sync"
	"syscall/js"
	"time"
)

type game struct {
	mutex          sync.Mutex
	score          int
	player         *player
	army           *alienArmy
	aliveCount     int                   // live alien count
	activeMissiles map[*missile]struct{} // active alien missile
	startTime      time.Time
	isOver         bool
}

func NewGame() *game {
	return &game{
		player:         initialPlayer(),
		army:           initialArmy(),
		mutex:          sync.Mutex{},
		aliveCount:     armyRowCount * armyColumnCount,
		activeMissiles: map[*missile]struct{}{},
		startTime:      time.Now(),
	}
}

func (g *game) listenForArrowKeys() {
	for {
		if g.isOver {
			break
		}
		g.mutex.Lock()
		switch {
		case keyDownState[_LEFT]:
			if g.player.baseY-1 > 20 {
				g.player.clear()
				g.player.baseY -= 1
				g.player.render()
				// g.player.move(_LEFT)
			}
		case keyDownState[_UP]:
			if g.player.baseX-1 > canvasDimX-30 {
				g.player.clear()
				g.player.baseX -= 1
				g.player.render()
				// g.player.move(_UP)
			}
		case keyDownState[_RIGHT]:
			if g.player.baseY+1 < canvasDimY-20 {
				g.player.clear()
				g.player.baseY += 1
				g.player.render()
				// g.player.move(_RIGHT)
			}
		case keyDownState[_DOWN]:
			if g.player.baseX+1 < canvasDimX-10 {
				g.player.clear()
				g.player.baseX += 1
				g.player.render()
				// g.player.move(_DOWN)
			}
		}
		g.mutex.Unlock()
		time.Sleep(time.Millisecond * 40)
	}
}

func (g *game) listenForSpaceKey() {
	for {
		if g.isOver {
			break
		}
		if keyDownState[5] {
			go func(m *missile) { // player shootMissile
				m.render()
				for {
					g.mutex.Lock()
					m.clear()
					m.move()
					m.destroyed = g.checkMissileCollision(m)
					if m.destroyed || m.baseX < 0 || g.checkMissileHit(m) || g.isOver {
						m.destroyed = true
						g.mutex.Unlock()
						break
					}
					m.render()
					g.mutex.Unlock()
					time.Sleep(time.Millisecond * 5)
				}
			}(g.player.launchMissile())
			time.Sleep(time.Millisecond * 490)
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func (g *game) checkMissileCollision(playerMissile *missile) bool {
	for m := range g.activeMissiles {
		if m.destroyed {
			continue
		}
		for _, c1 := range m.getBaseCoordinates() {
			for _, c2 := range playerMissile.getBaseCoordinates() {
				if c1.add(m.baseX, m.baseY) == c2.add(playerMissile.baseX, playerMissile.baseY) {
					m.destroyed = true
					playerMissile.destroyed = true
					return true
				}
			}
		}
	}
	return false
}

func (g *game) checkMissileHit(m *missile) bool {
	if m.kind == playerMissile {
		for idx, a := range g.army.aliens {
			if !a.died {
				for _, c := range a.getBaseCoordinates() {
					for _, c2 := range m.getBaseCoordinates() {
						if newCoord(m.baseX+c2.x, m.baseY+c2.y) == newCoord(c.x+a.baseX, c.y+a.baseY) {
							a.clear(a.state)
							g.army.aliens[idx].died = true
							g.aliveCount -= 1
							g.score += alienPoints[a.kind]
							g.renderScore()
							return true
						}
					}
				}
			}
		}
	} else {
		for _, c := range g.player.getBaseCoordinates() {
			for _, c2 := range m.getBaseCoordinates() {
				if c.add(g.player.baseX, g.player.baseY) == c2.add(m.baseX, m.baseY) {
					g.isOver = true
					return true
				}
			}
		}
	}
	return false
}

func (g *game) activateAlienArmy() {
	for id := range g.army.aliens {
		go g.activateAlienMissile(id)
	}
}
func (g *game) activateAlienMissile(id /*alien id*/ int) {
	for { // shoot missile
		if g.isOver {
			break
		}
		t := g.army.getLaunchDelay() - g.launchTimeDiff()
		hop := time.Millisecond * 50
		for {
			if t < hop {
				time.Sleep(t)
				break
			} else {
				t -= hop
				time.Sleep(hop)
			}
			if g.isOver {
				break
			}
		}
		if g.army.aliens[id].died {
			time.Sleep(time.Millisecond * 50)
			continue
		}
		go g.handleMissile(id)
	}
}

func (g *game) handleMissile(id int) {
	g.mutex.Lock()
	m := g.army.aliens[id].launchMissile()
	g.activeMissiles[m] = struct{}{}
	g.mutex.Unlock()
	m.render()
	for {
		g.mutex.Lock()
		m.clear()
		m.move()
		if m.destroyed || m.baseX >= canvasDimX-3 || g.isOver || g.checkMissileHit(m) {
			m.destroyed = true
			delete(g.activeMissiles, m)
			g.mutex.Unlock()
			break
		}
		m.render()
		g.mutex.Unlock()
		time.Sleep(time.Millisecond * 30)
	}
}

func (g *game) alienAtack() {
	for {
		g.mutex.Lock()
		g.army.render()
		g.mutex.Unlock()
		time.Sleep(time.Millisecond * 300)
		if g.isOver {
			break
		}
		g.mutex.Lock()
		g.army.clear()
		g.army.move(1)
		g.mutex.Unlock()
	}
}

// milliseconds time
func (g *game) launchTimeDiff() time.Duration {
	return time.Duration(time.Now().Sub(g.startTime).Seconds()) * time.Millisecond
}

func (g *game) initilize() {
	g.activateAlienArmy()
	go g.alienAtack()
	go g.listenForArrowKeys()
	go g.listenForSpaceKey()
	g.mutex.Lock()
	g.player.render()
	g.mutex.Unlock()
	g.init()
}

func (g *game) resetAlienArmy() {
	g.army = initialArmy()
	g.aliveCount = armyRowCount * armyColumnCount
}

func (g *game) Start() {
	g.initilize()
	g.renderScore()
	for {
		select {
		default:
			time.Sleep(time.Millisecond * 100)
			if g.isOver {
				g.mutex.Lock()
				g.clear()
				g.renderGameOver()
				g.mutex.Unlock()

				time.Sleep(time.Millisecond * 400)

				g.clear()
				g = NewGame()
				g.initilize()
				g.renderScore()
			}
			if g.aliveCount == 0 {
				time.Sleep(time.Millisecond * 200)
				g.mutex.Lock()
				g.resetAlienArmy()
				g.army.render()
				g.mutex.Unlock()
			}
		}
	}
}

var (
	_INIT        = false
	keyDownState = []bool{
		_LEFT:  false,
		_UP:    false,
		_RIGHT: false,
		_DOWN:  false,
		5:      false, // space key
	}
)

func (g *game) renderGameOver() {
	gameCanvas.Set("fillStyle", "green")
	gameCanvas.Set("font", "40px Arial")
	gameCanvas.Call("fillText", "Game Over", scale*canvasDimY/2, scale*canvasDimX/2-scale*20)
	gameCanvas.Call("fill")
}

func (g *game) renderScore() {
	gameCanvas.Call("clearRect", 0, 0, canvasDimY*scale, 15*scale)
	gameCanvas.Set("fillStyle", "green")
	gameCanvas.Set("font", "20px Arial")
	gameCanvas.Call("fillText", fmt.Sprintf("Score: %v", g.score), 2*scale, 6*scale)
	gameCanvas.Call("fill")
}

func (g *game) clear() {
	for x := range canvasDimX {
		for y := range canvasDimY {
			gameCanvas.Call("clearRect", y*scale, x*scale, scale, scale)
			missileCanvas.Call("clearRect", y*scale, x*scale, scale, scale)
		}
	}
	gameCanvas.Call("fill")
	missileCanvas.Call("fill")
}

func (g *game) init() {
	if _INIT {
		return
	}
	_INIT = true
	js.Global().Get("document").Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) any {
		keyCode := args[0].Get("keyCode").Int()
		g.mutex.Lock()
		switch keyCode {
		case 37, 38, 39, 40:
			keyDownState[keyCode-37+1] = true
		case 32:
			keyDownState[5] = true
		}
		g.mutex.Unlock()
		return nil
	}))

	js.Global().Get("document").Call("addEventListener", "keyup", js.FuncOf(func(this js.Value, args []js.Value) any {
		keyCode := args[0].Get("keyCode").Int()
		g.mutex.Lock()
		switch keyCode {
		case 37, 38, 39, 40:
			keyDownState[keyCode-37+1] = false
		case 32:
			keyDownState[5] = false
		}
		g.mutex.Unlock()
		return nil
	}))
}
