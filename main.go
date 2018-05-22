package main

import (
	"fmt"
	"net/http"
	"os"

	bencode "github.com/jackpal/bencode-go"
)

func main() {
	// read path to .torrent file from command line
	torrentPath := os.Args[1]

	torrentFile, err := os.Open(torrentPath)
	if err != nil {
		panic(err)
	}

	// decode the torrent file using BEencode
	m, err := bencode.Decode(torrentFile)
	if err != nil {
		panic(err)
	}

	metadata := m.(map[string]interface{})

	// extract useful things from the metadata
	announceURL := metadata["announce"].(string)
	// info := metadata["info"].(string)

	client := http.DefaultClient
	req, _ := http.NewRequest(http.MethodGet, announceURL, nil)

	// make GET request to tracker
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.StatusCode)
}
