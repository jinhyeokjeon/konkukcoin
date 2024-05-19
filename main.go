package main

import (
	"fmt"
	"os"

	"github.com/jinhyeokjeon/konkukcoin/p2p"
	"github.com/jinhyeokjeon/konkukcoin/ui"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <address>")
		os.Exit(1)
	}
	p2p.Port = os.Args[1]
	go p2p.Start()
	ui.Start()
}
