package loader

import (
	"errors"
	"os"

	bencode "github.com/jackpal/bencode-go"
)

type TorrentLoader interface {
	Load(string) error
	Decode() (map[string]interface{}, error)
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
		return err
	}

	f.file = torrentFile
	return nil
}

func (f *TorrentLoaderImpl) Decode() (map[string]interface{}, error) {
	m, err := bencode.Decode(f.file)
	if err != nil {
		return nil, err
	}

	metadataMap, ok := m.(map[string]interface{})
	if !ok {
		return nil, errors.New("Metadata couldn't be read as a map - likely invalid")
	}

	return metadataMap, nil
}
