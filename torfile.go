package btgo

import (
	"bytes"
	"errors"
	"math/rand"
	"os"
)

type File struct {
	path   string
	length int
}

type Torfile struct {
	files        []File
	announceList [][]string
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

	var announceList [][]string
	if announceListList := m["announce-list"]; announceListList != nil {
		if announceInterfaceListList, ok := announceListList.([]interface{}); ok {
			announceList = make([][]string, len(announceInterfaceListList))
			for i, e := range announceInterfaceListList {
				if announceInterfaceList, ok := e.([]interface{}); ok {
					innerAnnounceList := make([]string, len(announceInterfaceList))
					for in, el := range announceInterfaceList {
						if announceBytes, ok := el.([]byte); ok {
							announce := string(announceBytes)
							innerAnnounceList[in] = announce
						} else {
							err = errors.New("Unable to parse announce bytes")
							return
						}
					}
					shuffleStrings(innerAnnounceList)
					announceList[i] = innerAnnounceList
				} else {
					err = errors.New("Unable to parse inner announce list")
					return
				}
			}
		} else {
			err = errors.New("Unable to parse outer announce list")
			return
		}
	} else {
		announce, ok := stringFromBytesInterface(m["announce"])
		if !ok {
			err = errors.New("Unable to parse announce section of torfile")
			return
		}
		announceList = [][]string{[]string{announce}}
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
	name, ok := stringFromBytesInterface(info["name"])
	if !ok {
		err = errors.New("Unable to parse path for single-file torrent")
		return
	}

	if info["files"] == nil {
		length, ok := info["length"].(int)
		if !ok {
			err = errors.New("Unable to parse length for single-file torrent")
			return
		}
		files = []File{File{name, length}}
	} else {
		filesInterfaceList, ok := info["files"].([]interface{})
		if !ok {
			err = errors.New("Unable to parse files for multiple-file torrent")
			return
		}

		files = make([]File, len(filesInterfaceList))
		for i, e := range filesInterfaceList {
			fileInfo, ok := e.(map[string]interface{})
			if !ok {
				err = errors.New("Unable to parse file info for multiple-file torrent")
				return
			}

			length, ok := fileInfo["length"].(int)
			if !ok {
				err = errors.New("Unable to parse file lenfth in multiple-file torrent")
			}

			pathInterfaces, ok := fileInfo["path"].([]interface{})
			if !ok {
				err = errors.New("Unable to parse file path for multiple-file torrent")
				return
			}
			var pathBuffer bytes.Buffer
			pathBuffer.WriteString(name)
			pathBuffer.WriteRune(os.PathSeparator)
			for in, el := range pathInterfaces {
				pathPiece, ok := el.([]byte)
				if !ok {
					err = errors.New("Unable to parse file path piece for multiple-file torrent")
					return
				}
				pathBuffer.Write(pathPiece)
				if in != len(pathInterfaces)-1 {
					pathBuffer.WriteRune(os.PathSeparator)
				}
			}

			files[i] = File{pathBuffer.String(), length}
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

func shuffleStrings(strings []string) {
	for i := range strings {
		j := rand.Intn(i + 1)
		strings[i], strings[j] = strings[j], strings[i]
	}
}
