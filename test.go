package main

import (
	"EasyVoting/gui"
	"EasyVoting/ipfs"
)

//sudo sysctl -w net.core.rmem_max=2500000
func main() {
	is := ipfs.New(".ipfs")
	g := gui.New("EasyVoting", is)
	defer g.Close()

	g.Run()
}
