package btgo

import (
	"fmt"
	"reflect"
	"strconv"
)

func Buncode(s []byte) (buncoded interface{}) {
	buncoded, _ = buncode(s, 0)
	return
}

func buncode(st []byte, begin int) (buncoded interface{}, consumed int) {
	elements := make([]interface{}, 0)

OUTER:
	for i := begin; i < len(st); {
		switch {
		default:
			panic(fmt.Sprintf("Didn't expect character '%c' at position %d", st[i], i))
		case st[i] == 'i':
			val, n := parseInt(st, i)
			elements = append(elements, val)
			i, consumed = i+n, consumed+n
		case st[i] >= '0' && st[i] <= '9':
			str, n := parseString(st, i)
			elements = append(elements, str)
			i, consumed = i+n, consumed+n
		case st[i] == 'l':
			inner, n := buncode(st, i+1)
			elements = append(elements, inner)
			i, consumed = i+n, consumed+n
		case st[i] == 'd':
			innerElements, n := buncode(st, i+1)
			v := reflect.ValueOf(innerElements)
			d := make(map[string]interface{})
			for j := 0; j < v.Len(); j += 2 {
				if s, ok := v.Index(j).Interface().([]byte); ok {
					d[string(s)] = v.Index(j + 1).Interface()
				}
			}
			elements = append(elements, d)
			i, consumed = i+n, consumed+n
		case st[i] == 'e':
			i, consumed = i+2, consumed+2
			break OUTER
		}
	}

	if len(elements) == 1 {
		buncoded = elements[0]
	} else {
		buncoded = elements
	}
	return
}

func parseInt(st []byte, begin int) (val int, consumed int) {
	end := begin
	for st[end] != 'e' {
		end++
	}
	val, _ = strconv.Atoi(string(st[begin+1 : end]))
	consumed = end - begin + 1
	return
}

func parseString(st []byte, begin int) (s []byte, consumed int) {
	end := begin
	for st[end] >= '0' && st[end] <= '9' {
		end++
	}
	lengthStr := string(st[begin:end])
	length, _ := strconv.Atoi(lengthStr)
	end += 1 + length
	strBegin := begin + len(lengthStr) + 1
	s = st[strBegin : strBegin+length]
	consumed = end - begin
	return
}
