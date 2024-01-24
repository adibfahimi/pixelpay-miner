package main

import (
	"flag"
	"fmt"

	"github.com/adibfahimi/pixelpay-miner/core"
	"github.com/adibfahimi/pixelpay-miner/mine"
)

func main() {
	isSolo := flag.Bool("solo", false, "Mine in solo mode")
	nodeAddress := flag.String("node-address", "", "Set node address for solo mining")
	isPool := flag.Bool("pool", false, "Mine in pool mode")
	poolAddress := flag.String("pool-address", "", "Set pool address for pool mining")
	isWallet := flag.Bool("wallet", false, "Show your wallet")
	isAddress := flag.String("address", "", "Set wallet address")

	flag.Parse()

	if *isSolo {
		mine.MineSolo(*nodeAddress)
	} else if *isPool {
		mine.MinePool(*poolAddress)
	} else if *isWallet {
		core.ShowWallet()
	} else if *isAddress != "" {
		core.LoadWallet()
	} else {
		fmt.Println("no mode selected. For more information, use --help")
	}
}
