package main

import (
	"encoding/json"
	. "fmt"
	"github.com/algorand/go-algorand-sdk/client/algod"
	"os"
)

type Config struct {
	Token string `yaml:"token" action:"prompt"`
	Host  string `yaml:"host" action:"prompt,url"`
}

func main() {

	config := Config{}

	if err := LoadConfig(".algoc", &config); IsConfigNotPresent(err) {
		if err := PromptForValues(&config); err != nil {
			panic(err)
		}
		if err := WriteConfig(".algoc", config); err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	Printf("Interacting with network located at %s\n", config.Host)
	Println()

	var headers []*algod.Header
	headers = append(headers, &algod.Header{"X-API-Key", config.Token})
	// Create an algod client
	algodClient, err := algod.MakeClientWithHeaders(config.Host, "", headers)
	if err != nil {
		Printf("failed to make algod client: %s\n", err)
		return
	}

	// Print algod status
	Println("Status")
	nodeStatus, err := algodClient.Status()
	if err != nil {
		Printf("error getting algod status: %s\n", err)
		os.Exit(1)
	}

	Printf("algod last round: %d\n", nodeStatus.LastRound)
	Printf("algod time since last round: %d\n", nodeStatus.TimeSinceLastRound)
	Printf("algod catchup: %d\n", nodeStatus.CatchupTime)
	Printf("algod latest version: %s\n", nodeStatus.LastVersion)

	// Fetch block information
	lastBlock, err := algodClient.Block(nodeStatus.LastRound)
	if err != nil {
		Printf("error getting last block: %s\n", err)
		os.Exit(1)
	}

	// Print the block information
	Printf("\n-----------------Block Information-------------------\n")
	blockJSON, err := json.MarshalIndent(lastBlock, "", "\t")
	if err != nil {
		Printf("Can not marshall block data: %s\n", err)
	}
	Printf("%s\n", blockJSON)
}
