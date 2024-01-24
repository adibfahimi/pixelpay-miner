package mine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/adibfahimi/pixelpay-miner/core"
	pcore "github.com/adibfahimi/pixelpay-node/core"
)

type getBlockResponse struct {
	Message string       `json:"message"`
	Data    getBlockData `json:"data"`
}

type getBlockData struct {
	Block      pcore.Block `json:"block"`
	MineReward uint        `json:"reward"`
}

func getBlock(nodeAddress string) getBlockResponse {
	url := fmt.Sprintf("%s/mine", nodeAddress)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalf("client: could not create request: %s\n", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("client: could not send request: %s\n", err)
	}

	defer res.Body.Close()

	var data getBlockResponse

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		log.Fatalf("client: could not decode response: %s\n", err)
	}

	return data
}

func sendBlock(block pcore.Block, nodeAddress string) error {
	json, err := json.Marshal(block)
	if err != nil {
		log.Println(err)
	}

	url := fmt.Sprintf("%s/mine", nodeAddress)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return fmt.Errorf("client: could not create request: %s", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("client: could not send request: %s", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("client: unexpected status code: %d", res.StatusCode)
	}

	return nil
}

func MineSolo(nodeAddress string) {
	log.Println("Mining block ...")
	w := core.LoadWallet()

	for {
		log.Println("Getting block ...")
		d := getBlock(nodeAddress)
		log.Println("Block received.")

		if d.Message == "no pending tx" {
			log.Println("No pending tx. Waiting for 15 seconds ...")
			time.Sleep(15 * time.Second)
		} else {
			coinBaseTx := pcore.Tx{
				From:      "",
				To:        w.Address,
				Signature: "",
				Hash:      "",
				Amount:    d.Data.MineReward,
				Timestamp: uint(time.Now().Unix()),
			}

			coinBaseTx.Hash = coinBaseTx.CalculateHash()
			d.Data.Block.Txs = append(d.Data.Block.Txs, coinBaseTx)
			d.Data.Block.MerkleRoot = d.Data.Block.CalculateMerkleRoot()

			log.Println("Mining block ...")
			d.Data.Block.MineBlock()
			log.Println("Block mined. With nonce:", d.Data.Block.Nonce)

			if err := sendBlock(d.Data.Block, nodeAddress); err != nil {
				log.Fatalf("client: could not send block: %s\n", err)
			}

			log.Println("Block sent. mining next block ...")
		}
	}
}
