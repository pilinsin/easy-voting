package main

import (
	"EasyVoting/gui"
)

//sudo sysctl -w net.core.rmem_max=2500000
func main() {
	g := gui.New("EasyVoting", 810, 520)
	defer g.Close()

	g.Run()
}
