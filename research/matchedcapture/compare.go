package matchedcapture

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"unicode/utf8"
)

// Comparison covers the complete JSON value; descriptor strings remain opaque.
type Comparison struct {
	ByteEqual, StructuralEqual bool
	Paths                      []string
}

// CompareJSON ignores object order and whitespace, but never drops fields,
// sorts arrays, coerces numbers through floating point, or parses strings.
func CompareJSON(a, b []byte) (Comparison, error) {
	r := Comparison{ByteEqual: bytes.Equal(a, b), Paths: []string{}}
	x, e := parseJSON(a)
	if e != nil {
		return r, e
	}
	y, e := parseJSON(b)
	if e != nil {
		return r, e
	}
	differences(x, y, "$", &r.Paths)
	r.StructuralEqual = len(r.Paths) == 0
	return r, nil
}
func parseJSON(b []byte) (any, error) {
	if !utf8.Valid(b) {
		return nil, fmt.Errorf("invalid UTF-8")
	}
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	v, e := jsonValue(d)
	if e != nil {
		return nil, e
	}
	if _, e = d.Token(); e != io.EOF {
		return nil, fmt.Errorf("trailing JSON: %v", e)
	}
	return v, nil
}
func jsonValue(d *json.Decoder) (any, error) {
	t, e := d.Token()
	if e != nil {
		return nil, e
	}
	switch t {
	case json.Delim('{'):
		m := map[string]any{}
		for d.More() {
			k, e := d.Token()
			if e != nil {
				return nil, e
			}
			s, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("object key")
			}
			if _, ok = m[s]; ok {
				return nil, fmt.Errorf("duplicate key %q", s)
			}
			v, e := jsonValue(d)
			if e != nil {
				return nil, e
			}
			m[s] = v
		}
		_, e = d.Token()
		return m, e
	case json.Delim('['):
		a := []any{}
		for d.More() {
			v, e := jsonValue(d)
			if e != nil {
				return nil, e
			}
			a = append(a, v)
		}
		_, e = d.Token()
		return a, e
	default:
		return t, nil
	}
}
func differences(a, b any, path string, out *[]string) {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		*out = append(*out, path)
		return
	}
	switch x := a.(type) {
	case map[string]any:
		y := b.(map[string]any)
		keys := map[string]bool{}
		for k := range x {
			keys[k] = true
		}
		for k := range y {
			keys[k] = true
		}
		ordered := []string{}
		for k := range keys {
			ordered = append(ordered, k)
		}
		sort.Strings(ordered)
		for _, k := range ordered {
			av, ao := x[k]
			bv, bo := y[k]
			p := memberPath(path, k)
			if !ao || !bo {
				*out = append(*out, p)
			} else {
				differences(av, bv, p, out)
			}
		}
	case []any:
		y := b.([]any)
		for i := 0; i < max(len(x), len(y)); i++ {
			p := path + "[" + strconv.Itoa(i) + "]"
			if i >= len(x) || i >= len(y) {
				*out = append(*out, p)
			} else {
				differences(x[i], y[i], p, out)
			}
		}
	default:
		if !reflect.DeepEqual(a, b) {
			*out = append(*out, path)
		}
	}
}

func memberPath(parent, key string) string {
	identifier := key != ""
	for i, r := range key {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == '_' || i > 0 && r >= '0' && r <= '9') {
			identifier = false
			break
		}
	}
	if identifier {
		return parent + "." + key
	}
	return parent + "[" + string(jsonBytes(key)) + "]"
}
