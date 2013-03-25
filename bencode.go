package btgo

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func Bencode(t interface{}) (bencoded string) {
	switch k := reflect.TypeOf(t).Kind(); k {
	default:
		panic(fmt.Sprintf("unexpected type passed to Bencode: %T", t))
	case reflect.String:
		if st, ok := t.(string); ok {
			bencoded = bencodeString(st)
		}
	case reflect.Int:
		if st, ok := t.(int); ok {
			bencoded = bencodeInt(st)
		}
	case reflect.Slice:
		v := reflect.ValueOf(t)
		st := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			st[i] = interface{}(v.Index(i).Interface())
		}
		bencoded = bencodeSlice(st)
	case reflect.Map:
		v := reflect.ValueOf(t)
		st := make(map[string]interface{})
		keys := v.MapKeys()
		for _, key := range keys {
			st[key.String()] = v.MapIndex(key).Interface()
		}
		bencoded = bencodeMap(st)
	}

	return
}

func bencodeString(s string) string {
	return fmt.Sprintf("%d:%s", len(s), s)
}

func bencodeInt(i int) string {
	return fmt.Sprintf("i%de", i)
}

func bencodeSlice(s []interface{}) string {
	parts := make([]string, len(s))
	for i, e := range s {
		parts[i] = Bencode(e)
	}
	joined := strings.Join(parts, "")
	return fmt.Sprintf("l%se", joined)
}

func bencodeMap(m map[string]interface{}) string {
	parts := make([]string, len(m)*2)
	keys := keysInMap(m)
	sort.Strings(keys)
	i := 0
	for _, key := range keys {
		parts[i] = Bencode(key)
		parts[i+1] = Bencode(m[key])
		i += 2
	}
	joined := strings.Join(parts, "")
	return fmt.Sprintf("d%se", joined)
}

func keysInMap(m map[string]interface{}) []string {
	keys := make([]string, len(m))
	i := 0
	for key, _ := range m {
		keys[i] = key
		i++
	}
	return keys
}
