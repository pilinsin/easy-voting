package main

import (
	"EasyVoting/gui"
	"EasyVoting/ipfs"
	"EasyVoting/util"
)

//sudo sysctl -w net.core.rmem_max=2500000
func main() {
	is, err := ipfs.New(".ipfs")
	util.CheckError(err)
	g := gui.New("EasyVoting", is)
	defer g.Close()

	g.Run()
}
