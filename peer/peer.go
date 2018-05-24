package peer

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/go-errors/errors"
	"github.com/henryaj/gotorrent/tracker"
)

type PeerClient interface {
	ConnectAndDownload() error
}

type PeerClientImpl struct {
	file     *os.File
	peers    []tracker.Peer
	peerID   []byte
	metadata map[string]interface{}
	infoHash []byte
}

func NewPeerClient(file *os.File, peers []tracker.Peer, peerID []byte, metadata map[string]interface{}, infoHash []byte) PeerClient {
	return &PeerClientImpl{
		file:     file,
		peers:    peers,
		peerID:   peerID,
		metadata: metadata,
		infoHash: infoHash,
	}
}

func (c *PeerClientImpl) ConnectAndDownload() error {
	if len(c.peers) == 0 {
		return errors.New("no peers available")
	}

	for _, peer := range c.peers {
		fmt.Println(peer.IP, peer.Port)
		err := handshakeWithPeer(peer.IP, peer.Port, c.infoHash, c.peerID)
		if err != nil {
			return errors.Wrap(err, 0)
		}
	}
	return nil
}

func handshakeWithPeer(address string, port int, infoHash, peerID []byte) error {
	addr := fmt.Sprintf("[" + address + "]:" + strconv.Itoa(port))
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	handshake := constructHandshake(infoHash, peerID)

	defer conn.Close()
	binary.Write(conn, binary.BigEndian, &handshake)

	if err != nil {
		return errors.Wrap(err, 0)
	}

	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)

	log.Printf("Receive: %s", buf[:n])

	return nil

}

type handshake struct {
	len      uint8
	protocol []byte
	reserved [8]uint8
	infoHash [20]byte
	peerID   [20]byte
}

func constructHandshake(infoHash, peerID []byte) *handshake {
	protoIdentifier := []byte("BitTorrent protocol")

	h := &handshake{
		len:      uint8(len(protoIdentifier)),
		protocol: protoIdentifier,
		reserved: [8]uint8{0, 0, 0, 0, 0, 0, 0, 0},
	}

	copy(h.infoHash[:], infoHash)
	copy(h.peerID[:], peerID)

	return h
}
