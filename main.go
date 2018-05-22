package main

import (
	"os"

	"github.com/henryaj/gotorrent/loader"
	"github.com/henryaj/gotorrent/tracker"
)

func main() {
	torrentPath := os.Args[1]

	l := loader.NewTorrentLoader()
	err := l.Load(torrentPath)
	if err != nil {
		panic(err)
	}

	metadata, err := l.Decode()
	if err != nil {
		panic(err)
	}

	t := tracker.NewTrackerClient(metadata)

	_, err = t.GetPeers()
	if err != nil {
		panic(err)
	}
}
