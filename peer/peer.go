package peer

import (
	"fmt"
	"io"
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
	// addr, err := net.ResolveTCPAddr("tcp", address)
	// if err != nil {
	// 	return errors.Wrap(err, 0)
	// }

	addr := fmt.Sprintf("[" + address + "]:" + strconv.Itoa(port))
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	handshake := constructHandshake(infoHash, peerID)

	_, err = conn.Write(handshake)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	var response []byte
	for {
		_, err := conn.Read(response)
		if err != nil {
			if err != io.EOF {
				return errors.Wrap(err, 0)
			}
			break
		}
	}

	return nil

}

func constructHandshake(infoHash, peerID []byte) []byte {
	protoIdentifier := []byte("BitTorrent protocol")
	identifierLength := rune(len(protoIdentifier))
	reserved := []byte{0, 0, 0, 0, 0, 0, 0, 0}

	var handshake []byte
	handshake = append(handshake, byte(identifierLength))
	handshake = append(handshake, protoIdentifier...)
	handshake = append(handshake, reserved...)
	handshake = append(handshake, infoHash...)
	handshake = append(handshake, peerID...)

	return handshake
}
