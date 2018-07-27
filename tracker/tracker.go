package tracker

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/go-errors/errors"
	bencode "github.com/jackpal/bencode-go"
)

type TrackerClient interface {
	GetPeers() ([]string, error)
}

type TrackerClientImpl struct {
	metadata map[string]interface{}
	clientID []byte
}

type announceResponse struct {
	Complete   int    `bencode:"complete"`
	Incomplete int    `bencode:"incomplete"`
	Interval   int    `bencode:"interval"`
	Peers      []Peer `bencode:"peers"`
}

type Peer struct {
	IP         string `bencode:"complete"`
	PeerID     []byte `bencode:"peer id"`
	Port       int    `bencode:"port"`
	Connection *net.Conn
}

func NewTrackerClient(m map[string]interface{}, clientID []byte) *TrackerClientImpl {
	return &TrackerClientImpl{
		metadata: m,
		clientID: clientID,
	}
}

func (t *TrackerClientImpl) GetPeers() ([]Peer, error) {
	infoBlock := t.metadata["info"].(map[string]interface{})

	hash, err := calculateInfoHash(infoBlock)

	if err != nil {
		return nil, err
	}

	queryParams := map[string]string{
		"info_hash": string(hash),
		"peer_id":   string(t.clientID),
		"port":      "6881",
		"uploaded":  "0",
		"left":      "0",
		"compact":   "0",
		"event":     "started",
	}

	announceURL := t.metadata["announce"].(string)
	if announceURL == "" {
		return nil, errors.New("no announce URL found in torrent metadata")
	}

	fmt.Printf("Announce URL found: %s\n", announceURL)

	res, err := doRequest(announceURL, queryParams)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		bodyText, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(bodyText))
		errorMsg, err := bencode.Decode(res.Body)
		if err != nil {
			panic(errors.Wrap(err, 0))
		}
		return nil, fmt.Errorf("Error getting list of peers: %s\n", errorMsg)
	}

	announce := announceResponse{}
	err = bencode.Unmarshal(res.Body, &announce)
	if err != nil {
		panic(err)
	}

	prunedPeers := prunePeers(announce.Peers, t.clientID)

	return prunedPeers, nil
}

func prunePeers(peers []Peer, clientID []byte) []Peer {
	var prunedPeers []Peer
	for _, peer := range peers {
		// if bytes.Equal(peer.PeerID, clientID) {
		// 	continue
		// }
		if peer.Port == 6881 {
			continue // TODO: what the fuck
		}
		prunedPeers = append(prunedPeers, peer)
	}

	return prunedPeers
}

func calculateInfoHash(info map[string]interface{}) ([]byte, error) {
	encodedInfo := new(bytes.Buffer)
	bencode.Marshal(encodedInfo, info)

	hash := sha1.New()
	hash.Write(encodedInfo.Bytes())

	return hash.Sum(nil), nil
}

func doRequest(announceURL string, queryParams map[string]string) (*http.Response, error) {
	client := http.DefaultClient
	req, _ := http.NewRequest(http.MethodGet, announceURL, nil)
	q := req.URL.Query()

	for k, v := range queryParams {
		q.Set(k, v)
	}

	req.URL.RawQuery = q.Encode()

	return client.Do(req)
}
