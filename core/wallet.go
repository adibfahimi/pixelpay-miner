package core

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jaypipes/ghw"
)

type Wallet struct {
	Address string `json:"address"`
	Balance uint   `json:"balance"`
}

func LoadWallet() Wallet {
	file := "wallet.json"
	var wallet Wallet

	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println("wallet not found. Creating a new wallet...")
		fmt.Println("please enter your wallet address with the -address flag")
	} else {
		fmt.Println("wallet found. Loading wallet...")
		f, err := os.Open(file)
		if err != nil {
			log.Fatalf("an error occurred while opening file: %s", err)
		}
		defer f.Close()

		data := make([]byte, 100)
		count, err := f.Read(data)
		if err != nil {
			log.Fatalf("an error occurred while reading file: %s count: %d", err, count)
		}

		if err := json.Unmarshal(data[:count], &wallet); err != nil {
			log.Fatalf("an error occurred while unmarshalling JSON data: %s", err)
		}
	}

	return wallet
}

func (w *Wallet) SaveWallet() error {
	file := "wallet.json"
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("an error occurred while creating file: %s", err)
	}
	defer f.Close()

	data, err := json.Marshal(w)
	if err != nil {
		return fmt.Errorf("an error occurred while marshalling JSON data: %s", err)
	}

	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("an error occurred while writing to file: %s", err)
	}

	return nil
}

func ShowWallet() {
	w := LoadWallet()
	w.GetBalance()

	fmt.Println("your wallet address is:", w.Address)
	fmt.Println("your wallet balance is:", w.Balance)

	cpu, err := ghw.CPU()
	if err == nil {
		fmt.Println(cpu.String())
	}

	memory, err := ghw.Memory()
	if err == nil {
		fmt.Println(memory.String())
	}
}
