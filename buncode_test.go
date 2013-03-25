package btgo

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestBuncode(t *testing.T) {
	// Integers
	s := "i3e"
	if r := Buncode([]byte(s)); r != 3 {
		t.Errorf("Doesn't decode %s correctly: %v", s, r)
	}
	s = "i123456e"
	if r := Buncode([]byte(s)); r != 123456 {
		t.Errorf("Doesn't decode %s correctly: %v", s, r)
	}

	// Strings
	s = "1:a"
	if r := Buncode([]byte(s)); !sameSlice(r, []byte("a")) {
		t.Errorf("Doesn't decode %s correctly: %v", s, r)
	}
	s = "4:spam"
	if r := Buncode([]byte(s)); !sameSlice(r, []byte("spam")) {
		t.Errorf("Doesn't decode %s correctly: %v", s, r)
	}

	// Lists (slices)
	s = "l4:spam4:eggse"
	if r := Buncode([]byte(s)); !sameSlice(r, [][]byte{[]byte("spam"), []byte("eggs")}) {
		t.Errorf("Doesn't decode %s correctly: %v", s, r)
	}
	s = "l1:al3:wayi2e2:goe1:ce"
	if r := Buncode([]byte(s)); !sameSlice(r, []interface{}{[]byte("a"), []interface{}{[]byte("way"), 2, []byte("go")}, []byte("c")}) {
		t.Errorf("Doesn't decode %s correctly: %v", s, r)
	}

	// Dictionaries (maps)
	s = "d9:publisher3:bob17:publisher-webpage15:www.example.com18:publisher.location4:homee"
	d, ok := Buncode([]byte(s)).(map[string]interface{})
	if !ok && (!sameSlice(d["publisher"], []byte("bob")) || !sameSlice(d["publisher-webpage"], []byte("www.example.com")) || !sameSlice(d["publisher.location"], []byte("home"))) {
		t.Errorf("Doesn't decode %s correctly: %v", s, d)
	}

	// File tests (.torrent files)
	content, err := ioutil.ReadFile("test/ubuntu.torrent")
	worked := false
	if err == nil {
		r := Buncode(content)
		m, ok := r.(map[string]interface{})
		if ok {
			info, ok2 := m["info"].(map[string]interface{})
			if ok2 {
				name, ok3 := info["name"]
				length, ok4 := info["length"].(int)
				if ok3 && ok4 && sameSlice(name, []byte("ubuntu-12.10-desktop-amd64.iso")) && length == 800063488 {
					worked = true
				}
			}
		}
	}
	if !worked {
		t.Errorf("Doesn't decode ubuntu.torrent correctly")
	}
}

func sameSlice(a interface{}, b interface{}) bool {
	return (fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b))
}
