package main

import (
	"github.com/pilinsin/easy-voting/gui"
)

//sudo sysctl -w net.core.rmem_max=2500000
func main() {
	g := gui.New("EasyVoting", 450, 800)
	g.Run()
}
