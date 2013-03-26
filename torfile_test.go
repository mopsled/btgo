package btgo

import (
	"io/ioutil"
	"math"
	"testing"
)

func TestNewTorfile(t *testing.T) {
	// One-file torrent
	testUbuntuTorrent(t)

	// Multi-file torrent
	testMultitracksTorrent(t)
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

	if !sameSlice(tfile.announceList, [][]string{[]string{"http://torrent.ubuntu.com:6969/announce"}, []string{"http://ipv6.torrent.ubuntu.com:6969/announce"}}) {
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

func testMultitracksTorrent(t *testing.T) {
	file := "test/multitracks.torrent"
	content, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("Failed to open test file %s", file)
	}

	tfile, err := NewTorfile(content)
	if err != nil {
		t.Fatalf("Failed to parse test file %s: %s", file, err)
	}

	if !sameSlice(tfile.announceList, [][]string{[]string{"http://tracker001.clearbits.net:7070/announce"}}) {
		t.Errorf("Wrong announce list for %s: %s", file, tfile.announceList)
	}

	if tfile.pieceLength != 262144 {
		t.Errorf("Wrong pieceLength for %s: %d", file, tfile.pieceLength)
	}

	if len(tfile.files) != 6 {
		t.Fatalf("Wrong number of files for %s: %d", file, len(tfile.files))
	}
}
