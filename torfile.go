package btgo

import (
	"errors"
	"fmt"
)

type File struct {
	path   string
	length int
}

type Torfile struct {
	files        []File
	announceList []string
	pieceLength  int
	pieces       [][20]byte
}

func NewTorfile(file []byte) (tfile *Torfile, err error) {
	buncoded := Buncode(file)
	m, ok := buncoded.(map[string]interface{})
	if !ok {
		err = errors.New("Buncoding error: unable to parse torfile")
		return
	}

	announce, ok := stringFromBytesInterface(m["announce"])
	if !ok {
		err = errors.New("Unable to parse announce section of torfile")
		return
	}
	announceList := []string{announce}
	if announceListInterface := m["announce-list"]; announceListInterface != nil {
		if announceListBytes, ok := announceListInterface.([]interface{}); ok {
			for _, e := range announceListBytes {
				if es, ok := e.([]byte); ok {
					an := string(es)
					if !inStringSlice(an, announceList) {
						announceList = append(announceList, an)
					}
				}
			}
		}
	}

	info, ok := m["info"].(map[string]interface{})
	if !ok {
		err = errors.New("Unable to parse info dictionary in torfile")
		return
	}

	pieceLength, ok := info["piece length"].(int)
	if !ok {
		err = errors.New("Unable to parse piece length")
		return
	}

	pieceBytes, ok := info["pieces"].([]byte)
	if !ok {
		err = errors.New("Unable to parse piece hashes in torfile")
		return
	}
	hashes := len(pieceBytes) / 20
	pieces := make([][20]byte, hashes)
	for i := 0; i < hashes; i += 1 {
		var piece [20]byte
		copy(piece[:], pieceBytes[i*20:(i+1)*20])
	}

	files, err := filesFromInfo(info)
	if err != nil {
		return
	}

	tfile = &Torfile{files, announceList, pieceLength, pieces}
	return
}

func filesFromInfo(info map[string]interface{}) (files []File, err error) {
	if info["files"] == nil {
		path, ok := stringFromBytesInterface(info["name"])
		if !ok {
			err = errors.New("Unable to parse path for single-file torrent")
			return
		}

		length, ok := info["length"].(int)
		if !ok {
			err = errors.New("Unable to parse length for single-file torrent")
			return
		}
		files = []File{File{path, length}}
	} else {
		_, ok := info["files"].([]map[string]interface{})
		if ok {
			fmt.Println("files are okay!")
		}
	}

	return
}

func stringFromBytesInterface(i interface{}) (s string, ok bool) {
	var sBytes []byte
	if sBytes, ok = i.([]byte); ok {
		s = string(sBytes)
	}
	return
}

func inStringSlice(s string, sl []string) (inSlice bool) {
	for _, e := range sl {
		if e == s {
			inSlice = true
			return
		}
	}
	return
}
