package btgo

import (
	"io/ioutil"
	"math"
	"math/big"
	"testing"
)

func TestNewTorfile(t *testing.T) {
	// One-file torrents
	testUbuntuTorrent(t)

	// Multi-file torrents
	testMultitracksTorrent(t)
	testStackExchangeTorrent(t)
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

	expected := big.NewInt(524288)
	if tfile.pieceLength.Cmp(expected) != 0 {
		t.Errorf("Wrong pieceLength for %s: %d", file, tfile.pieceLength)
	}

	if len(tfile.files) != 1 {
		t.Fatalf("Wrong number of files for %s: %d", file, len(tfile.files))
	}
	firstFile := tfile.files[0]
	if firstFile.path != "ubuntu-12.10-desktop-amd64.iso" {
		t.Errorf("Wrong title for file in %s: %s", file, firstFile.path)
	}
	expected = big.NewInt(800063488)
	if firstFile.length.Cmp(expected) != 0 {
		t.Errorf("Wrong length for file in %s: %d", file, tfile.pieceLength)
	}

	if len(tfile.pieces) != int(math.Ceil(float64(firstFile.length.Int64())/float64(tfile.pieceLength.Int64()))) {
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

	expected := big.NewInt(262144)
	if tfile.pieceLength.Cmp(expected) != 0 {
		t.Errorf("Wrong pieceLength for %s: %d", file, tfile.pieceLength)
	}

	if len(tfile.files) != 6 {
		t.Fatalf("Wrong number of files for %s: %d", file, len(tfile.files))
	}

	expected = big.NewInt(753049736)
	if tfile.files[0].path != "Flembaz - Floppy Disk feat Stylver_multitracks/Content/Flembaz - 10 - Floppy Disk feat Stylver_multitracks/Content/10 - Floppy Disk (feat. Stylver) [Multitracks].rar" || tfile.files[0].length.Cmp(expected) != 0 {
		t.Errorf("Wrong file information for %s: (%s, %d)", file, tfile.files[0].path, tfile.files[0].length)
	}
	expected = big.NewInt(980)
	if tfile.files[1].path != "Flembaz - Floppy Disk feat Stylver_multitracks/Content/Flembaz - 10 - Floppy Disk feat Stylver_multitracks/Description.txt" || tfile.files[1].length.Cmp(expected) != 0 {
		t.Errorf("Wrong file information for %s: (%s, %d)", file, tfile.files[1].path, tfile.files[1].length)
	}
	expected = big.NewInt(45)
	if tfile.files[2].path != "Flembaz - Floppy Disk feat Stylver_multitracks/Content/Flembaz - 10 - Floppy Disk feat Stylver_multitracks/License.txt" || tfile.files[2].length.Cmp(expected) != 0 {
		t.Errorf("Wrong file information for %s: (%s, %d)", file, tfile.files[2].path, tfile.files[2].length)
	}
	expected = big.NewInt(57987)
	if tfile.files[3].path != "Flembaz - Floppy Disk feat Stylver_multitracks/Content/flembaz - 10 - floppy disk feat stylver_multitracks.torrent" || tfile.files[3].length.Cmp(expected) != 0 {
		t.Errorf("Wrong file information for %s: (%s, %d)", file, tfile.files[3].path, tfile.files[3].length)
	}
	expected = big.NewInt(1170)
	if tfile.files[4].path != "Flembaz - Floppy Disk feat Stylver_multitracks/Description.txt" || tfile.files[4].length.Cmp(expected) != 0 {
		t.Errorf("Wrong file information for %s: (%s, %d)", file, tfile.files[4].path, tfile.files[4].length)
	}
	expected = big.NewInt(45)
	if tfile.files[5].path != "Flembaz - Floppy Disk feat Stylver_multitracks/License.txt" || tfile.files[5].length.Cmp(expected) != 0 {
		t.Errorf("Wrong file information for %s: (%s, %d)", file, tfile.files[5].path, tfile.files[5].length)
	}
}

func testStackExchangeTorrent(t *testing.T) {
	file := "test/stack-exchange.torrent"
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

	expected := big.NewInt(8388608)
	if tfile.pieceLength.Cmp(expected) != 0 {
		t.Errorf("Wrong pieceLength for %s: %d", file, tfile.pieceLength)
	}

	if len(tfile.files) != 91 {
		t.Fatalf("Wrong number of files for %s: %d", file, len(tfile.files))
	}
	expected = big.NewInt(21122632)
	if tfile.files[0].path != "Stack Exchange Data Dump - Mar 2013/Content/android.stackexchange.com.7z" || tfile.files[0].length.Cmp(expected) != 0 {
		t.Errorf("Wrong file information for %s: (%s, %d)", file, tfile.files[0].path, tfile.files[0].length)
	}
	expected = big.NewInt(1089277)
	if tfile.files[50].path != "Stack Exchange Data Dump - Mar 2013/Content/meta.ux.stackexchange.com.7z" || tfile.files[50].length.Cmp(expected) != 0 {
		t.Errorf("Wrong file information for %s: (%s, %d)", file, tfile.files[50].path, tfile.files[50].length)
	}
	expected = big.NewInt(48)
	if tfile.files[90].path != "Stack Exchange Data Dump - Mar 2013/License.txt" || tfile.files[90].length.Cmp(expected) != 0 {
		t.Errorf("Wrong file information for %s: (%s, %d)", file, tfile.files[90].path, tfile.files[90].length)
	}
}
