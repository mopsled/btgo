package btgo

import (
	"io/ioutil"
	"math"
	"testing"
)

func TestNewTorfile(t *testing.T) {
	// One-file torrent
	testUbuntuTorrent(t)
}

func testUbuntuTorrent(t *testing.T) {
	file := "test/ubuntu.torrent"
	content, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("Failed to open test file %s", file)
	}

	tfile, err := NewTorfile(content)
	if err != nil {
		t.Fatalf("Failed to parse test file %s: %s", file, err)
	}

	if tfile.announce != "http://torrent.ubuntu.com:6969/announce" {
		t.Errorf("Wrong announce address for %s: %s", file, tfile.announce)
	}

	if !sameSlice(tfile.announceList, []string{"http://torrent.ubuntu.com:6969/announce", "http://ipv6.torrent.ubuntu.com:6969/announce"}) {
		t.Errorf("Wrong announce list for %s: %s", file, tfile.announceList)
	}

	if tfile.pieceLength != 524288 {
		t.Errorf("Wrong pieceLength for %s: %d", file, tfile.pieceLength)
	}

	if len(tfile.files) != 1 {
		t.Fatalf("Wrong number of files for %s: %d", file, len(tfile.files))
	}
	firstFile := tfile.files[0]
	if firstFile.path != "ubuntu-12.10-desktop-amd64.iso" {
		t.Errorf("Wrong title for file in %s: %s", file, firstFile.path)
	}
	if firstFile.length != 800063488 {
		t.Errorf("Wrong length for file in %s: %d", file, firstFile.length)
	}

	if len(tfile.pieces) != int(math.Ceil(float64(firstFile.length)/float64(tfile.pieceLength))) {
		t.Errorf("Wrong size of piece slice for %s: %d", file, len(tfile.pieces))
	}
}
