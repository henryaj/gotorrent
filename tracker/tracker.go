package tracker

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	bencode "github.com/jackpal/bencode-go"
)

type TrackerClient interface {
	GetPeers() ([]string, error)
}

type TrackerClientImpl struct {
	metadata map[string]interface{}
}

func NewTrackerClient(m map[string]interface{}) *TrackerClientImpl {
	return &TrackerClientImpl{
		metadata: m,
	}
}

func (t *TrackerClientImpl) GetPeers() ([]string, error) {
	infoBlock := t.metadata["info"].(map[string]interface{})

	hash, err := calculateInfoHash(infoBlock)

	if err != nil {
		return nil, err
	}

	queryParams := map[string]string{
		"info_hash": url.QueryEscape(string(hash)),
		"peer_id":   url.QueryEscape(string(newClientID())),
		"port":      "6881",
		"uploaded":  "0",
		"left":      "0",
		"compact":   "0",
		"event":     "started",
	}

	announceURL := t.metadata["announce"].(string)
	fmt.Println(announceURL)

	res, err := doRequest(announceURL, queryParams)

	if res.StatusCode != http.StatusAccepted {
		body, err := bencode.Decode(res.Body)
		if err != nil {
			body, _ := ioutil.ReadAll(res.Body)
			fmt.Println(string(body))
		}
		return nil, fmt.Errorf("Error getting list of peers: %s\n", body)
	}

	return nil, nil

}

func calculateInfoHash(info map[string]interface{}) ([]byte, error) {
	infoRaw := new(bytes.Buffer)

	err := bencode.Marshal(infoRaw, info)

	if err != nil {
		return nil, fmt.Errorf("Unable to bencode info dict: %s\n", err)
	}

	hasher := sha1.New()
	var infoHash []byte
	hasher.Write(infoRaw.Bytes())
	copy(infoHash, hasher.Sum(nil))

	return infoHash, nil
}

func doRequest(announceURL string, queryParams map[string]string) (*http.Response, error) {
	client := http.DefaultClient
	req, _ := http.NewRequest(http.MethodGet, announceURL+"/announce", nil)
	q := req.URL.Query()

	for k, v := range queryParams {
		q.Set(k, v)
	}

	req.URL.RawQuery = q.Encode()

	// make GET request to tracker
	return client.Do(req)
}

func newClientID() string {
	clientID := make([]byte, 20)
	rand.Read(clientID)
	return string(clientID)
}
