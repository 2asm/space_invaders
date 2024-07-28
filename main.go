//go:build js && wasm

package main

import "github.com/2asm/space_invaders/game"

func main() {
	game.NewGame().Start()
}
