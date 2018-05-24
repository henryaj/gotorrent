package main

import (
	"math/rand"
	"os"

	"github.com/go-errors/errors"
	"github.com/henryaj/gotorrent/loader"
	"github.com/henryaj/gotorrent/peer"
	"github.com/henryaj/gotorrent/tracker"
)

func main() {
	torrentPath := os.Args[1]

	l := loader.NewTorrentLoader()
	err := l.Load(torrentPath)
	if err != nil {
		panic(err)
	}

	metadata, infoHash, err := l.Decode()
	if err != nil {
		panic(err.(*errors.Error).ErrorStack())
	}

	clientID := newClientID()
	t := tracker.NewTrackerClient(metadata, clientID)

	peers, err := t.GetPeers()
	if err != nil {
		panic(err.(*errors.Error).ErrorStack())
	}

	pc := peer.NewPeerClient(nil, peers, clientID, metadata, infoHash)

	if err = pc.ConnectAndDownload(); err != nil {
		panic(err.(*errors.Error).ErrorStack())
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func newClientID() []byte {
	clientID := make([]rune, 20)
	for i := range clientID {
		clientID[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return []byte(string(clientID))
}
