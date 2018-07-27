package loader

import (
	"bytes"
	"crypto/sha1"

	"github.com/go-errors/errors"
	bencode "github.com/jackpal/bencode-go"
)

//TorrentLoader is a utility for decoding "bencoded" torrent metadata
type TorrentLoader interface {
	Decode([]byte) (map[string]interface{}, []byte, error)
}

type torrentLoaderImpl struct {
}

//NewTorrentLoader creates a new torrent loader
func NewTorrentLoader() TorrentLoader {
	return &torrentLoaderImpl{}
}

func (f *torrentLoaderImpl) Decode(inBytes []byte) (map[string]interface{}, []byte, error) {
	b := bytes.NewReader(inBytes)
	m, err := bencode.Decode(b)
	if err != nil {
		return nil, nil, errors.Wrap(err, 0)
	}

	metadataMap, ok := m.(map[string]interface{})
	if !ok {
		return nil, nil, errors.Wrap(errors.New("Metadata couldn't be read as a map - likely invalid"), 1)
	}

	hash, err := calculateInfoHash(metadataMap)
	if err != nil {
		return nil, nil, errors.Wrap(err, 0)
	}

	return metadataMap, hash, nil
}

func calculateInfoHash(info map[string]interface{}) ([]byte, error) {
	encodedInfo := new(bytes.Buffer)
	bencode.Marshal(encodedInfo, info)

	hash := sha1.New()
	hash.Write(encodedInfo.Bytes())

	return hash.Sum(nil), nil
}
