package btgo

import (
	"math/big"
	"testing"
)

func TestBencode(t *testing.T) {
	// Strings
	if s := Bencode("spam"); s != "4:spam" {
		t.Error("Doesn't encode spam right:", s)
	}
	if s := Bencode("eggs and such"); s != "13:eggs and such" {
		t.Error("Doesn't encode 'eggs and such' right:", s)
	}

	// Integers
	if s := Bencode(3); s != "i3e" {
		t.Error("Doesn't encode 3 correctly:", s)
	}
	if s := Bencode(100); s != "i100e" {
		t.Error("Doesn't encode 100 correctly:", s)
	}

	// Big Integers
	expected := new(big.Int)
	expected.SetString("123456789123456789", 10)
	if s := Bencode(expected); s != "i123456789123456789e" {
		t.Error("Doesn't encode 123456789123456789 correctly:", s)
	}

	// Lists (slices)
	stringSlice := []string{"spam", "eggs"}
	if s := Bencode(stringSlice); s != "l4:spam4:eggse" {
		t.Errorf("Doesn't encode %v correctly: %s", stringSlice, s)
	}
	singleElementSlice := genSlice("path", genSlice("file"))
	if s := Bencode(singleElementSlice); s != "l4:pathl4:fileee" {
		t.Errorf("Doesn't encode %v correctly: %s", stringSlice, s)
	}
	multiSlice := genSlice("way", 2, "go")
	if s := Bencode(multiSlice); s != "l3:wayi2e2:goe" {
		t.Errorf("Doesn't encode %v correctly: %s", multiSlice, s)
	}
	twoLayerSlice := genSlice("a", multiSlice, "c")
	if s := Bencode(twoLayerSlice); s != "l1:al3:wayi2e2:goe1:ce" {
		t.Errorf("Doesn't encode %v correctly: %s", twoLayerSlice, s)
	}

	// Dictionaries (maps)
	stringMap := map[string]string{"publisher": "bob", "publisher-webpage": "www.example.com", "publisher.location": "home"}
	if s := Bencode(stringMap); s != "d9:publisher3:bob17:publisher-webpage15:www.example.com18:publisher.location4:homee" {
		t.Errorf("Doesn't encode %v correctly: %s", stringMap, s)
	}
	listMap := map[string]interface{}{"spam": genSlice("a", "b")}
	if s := Bencode(listMap); s != "d4:spaml1:a1:bee" {
		t.Errorf("Doesn't encode %v correctly: %s", listMap, s)
	}

	// Mixup
	complexSlice := genSlice(genSlice("a", "b", "c"), listMap, genSlice(multiSlice, twoLayerSlice), "done")
	if s := Bencode(complexSlice); s != "ll1:a1:b1:ced4:spaml1:a1:beell3:wayi2e2:goel1:al3:wayi2e2:goe1:cee4:donee" {
		t.Errorf("Doesn't encode %v correctly: %s", listMap, s)
	}
}

func genSlice(args ...interface{}) (s []interface{}) {
	s = make([]interface{}, len(args))
	for i, e := range args {
		s[i] = interface{}(e)
	}
	return
}
