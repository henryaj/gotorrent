package loader

import (
	"bytes"
	"crypto/sha1"
	"os"

	"github.com/go-errors/errors"
	bencode "github.com/jackpal/bencode-go"
)

type TorrentLoader interface {
	Load(string) error
	Decode() (map[string]interface{}, []byte, error)
}

type TorrentLoaderImpl struct {
	path string
	file *os.File
}

func NewTorrentLoader() TorrentLoader {
	return &TorrentLoaderImpl{}
}

func (f *TorrentLoaderImpl) Load(path string) error {
	f.path = path
	torrentFile, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	f.file = torrentFile
	return nil
}

func (f *TorrentLoaderImpl) Decode() (map[string]interface{}, []byte, error) {
	m, err := bencode.Decode(f.file)
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
