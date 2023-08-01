package main

import (
	"Balloon/cfg"
	"Balloon/game"
	"math/rand"
	"runtime"
	"time"

	"github.com/wonderivan/logger"
)

func main() {
	rand.NewSource(time.Now().Unix())
	runtime.GOMAXPROCS(runtime.NumCPU())

	logger.SetLogger("logger.json")

	cfg := cfg.LoadCfg()
	if cfg == nil {
		return
	}

	game.Init(cfg)
}
