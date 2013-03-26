package btgo

import (
	"fmt"
	"math/big"
	"testing"
)

func TestBuncode(t *testing.T) {
	// Integers
	s := "i3e"
	expected := big.NewInt(3)
	r, ok := Buncode([]byte(s)).(*big.Int)
	if !ok || r.Cmp(expected) != 0 {
		t.Errorf("Doesn't decode %s correctly: %v", s, r)
	}
	s = "i123456e"
	expected = big.NewInt(123456)
	r, ok = Buncode([]byte(s)).(*big.Int)
	if !ok || r.Cmp(expected) != 0 {
		t.Errorf("Doesn't decode %s correctly: %v", s, r)
	}
	s = "i3306489856e"
	expected = new(big.Int)
	expected.SetString("3306489856", 10)
	r, ok = Buncode([]byte(s)).(*big.Int)
	if !ok || r.Cmp(expected) != 0 {
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
	if r := Buncode([]byte(s)); !sameSlice(r, []interface{}{[]byte("a"), []interface{}{[]byte("way"), big.NewInt(2), []byte("go")}, []byte("c")}) {
		t.Errorf("Doesn't decode %s correctly: %v", s, r)
	}
	s = "l4:pathl4:fileee"
	if r := Buncode([]byte(s)); !sameSlice(r, []interface{}{[]byte("path"), []interface{}{[]byte("file")}}) {
		t.Errorf("Doesn't decode %s correctly: %v", s, r)
	}
	s = "l4:look4:it's5:emptylee"
	if r := Buncode([]byte(s)); !sameSlice(r, []interface{}{[]byte("look"), []byte("it's"), []byte("empty"), []interface{}{}}) {
		t.Errorf("Doesn't decode %s correctly: %v", s, r)
	}

	// Dictionaries (maps)
	s = "d9:publisher3:bob17:publisher-webpage15:www.example.com18:publisher.location4:homee"
	d, ok := Buncode([]byte(s)).(map[string]interface{})
	if !ok && (!sameSlice(d["publisher"], []byte("bob")) || !sameSlice(d["publisher-webpage"], []byte("www.example.com")) || !sameSlice(d["publisher.location"], []byte("home"))) {
		t.Errorf("Doesn't decode %s correctly: %v", s, d)
	}
}

func sameSlice(a interface{}, b interface{}) bool {
	return (fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b))
}
