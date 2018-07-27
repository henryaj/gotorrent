package peer

import (
	"encoding/binary"
	"fmt"
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
		err := handshakeWithPeer(peer, c.infoHash, c.peerID)
		if err != nil {
			return errors.Wrap(err, 0)
		}
	}
	return nil
}

func handshakeWithPeer(peer tracker.Peer, infoHash, peerID []byte) error {
	addr := fmt.Sprintf("[" + peer.IP + "]:" + strconv.Itoa(peer.Port))
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	handshake := constructHandshake(infoHash, peerID)
	// make the handshake
	err = binary.Write(conn, binary.BigEndian, &handshake)

	// save the connection on the peer
	peer.Connection = &conn

	if err != nil {
		return errors.Wrap(err, 0)
	}

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
