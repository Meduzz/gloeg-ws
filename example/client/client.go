package main

import (
	"fmt"

	"github.com/Meduzz/gloegg"
	gloegws "github.com/Meduzz/gloegg-ws"
	"github.com/Meduzz/gloegg/toggles"
	"github.com/Meduzz/helper/block"
)

func main() {
	err := gloegws.Setup("", "test", "test")

	if err != nil {
		println("Gloegg WS setup failed", err.Error())
	}

	logger := gloegg.CreateLogger("Client")

	toggles.Subscribe(func(ut *toggles.UpdatedToggle) {
		fmt.Printf("toggle updated: %s\n", ut.Name)
	})

	logger.Info("Hello world!")
	logger.Error("A fake error occured", "error", fmt.Errorf("hi, im an error"))

	// keep app open to:
	// 1. send the error log
	// 2. allow us to play around with updating remote flags
	block.Block(func() error { return nil })
}
