package main

import (
	"flag"
	"fmt"

	"github.com/adibfahimi/pixelpay-miner/mine"
	"github.com/adibfahimi/pixelpay-miner/utils"
)

func main() {
	nodeAddress := flag.String("node", "", "set node address for solo mining")
	poolAddress := flag.String("pool", "", "set pool address for pool mining")
	address := flag.String("address", "", "set wallet address")

	flag.Parse()

	if *address == "" {
		fmt.Println("no wallet address provided")
		return
	}

	if !utils.IsValidAddress(*address) {
		fmt.Println("invalid wallet address")
		return
	}

	if *nodeAddress != "" {
		fmt.Println("start solo mining")
		mine.MineSolo(*nodeAddress, *address)
	} else if *poolAddress != "" {
		fmt.Println("start pool mining")
		mine.MinePool(*poolAddress, *address)
	} else {
		fmt.Println("no mode selected. For more information, use --help")
	}
}
