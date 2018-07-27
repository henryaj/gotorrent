package loader

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	tl := NewTorrentLoader()
	contents, err := ioutil.ReadFile("fixtures/alice.torrent")

	if err != nil {
		t.Error(err)
	}

	output, infoHash, err := tl.Decode(contents)

	info := output["info"].(map[string]interface{})
	if info["name"] != "alice.txt" {
		t.Error("failed to decode torrent")
	}

	if !reflect.DeepEqual(infoHash, []byte{105, 142, 104, 50, 143, 127, 31, 75, 208, 8, 112, 250, 108, 245, 172, 212, 183, 240, 237, 42}) {
		t.Error("incorrect info hash")
	}

	if err != nil {
		t.Error("error occurred: ", err)
	}
}
