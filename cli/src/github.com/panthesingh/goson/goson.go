package goson

import (
	"encoding/json"
	"strconv"
)

// Goson object
type Goson struct {
	i interface{}
}

/*
	bool, for JSON booleans
	float64, for JSON numbers
	string, for JSON strings
	[]interface{}, for JSON arrays
	map[string]interface{}, for JSON objects
	nil for JSON null
*/

func New(i interface{}) *Goson {
	return &Goson{i: i}
}

// Parse will create a goson object from json data
func Parse(data []byte) (*Goson, error) {
	var i interface{}
	if err := json.Unmarshal(data, &i); err != nil {
		return nil, err
	}
	return &Goson{i: i}, nil
}

// Get returns a goson object from a key.
// If the value does not exist this will still return a goson object.
func (g *Goson) Get(key string) *Goson {
	if m, ok := g.i.(map[string]interface{}); ok {
		return &Goson{i: m[key]}
	}
	return &Goson{i: new(interface{})}
}

// Set returns a goson object after change the value
// If the value does not exist this will still return a goson object.
func (g *Goson) Set(key string, value interface{}) *Goson {
	if m, ok := g.i.(map[string]interface{}); ok {
		m[key] = value
		return &Goson{i: m[key]}
	}
	return &Goson{i: new(interface{})}
}

// Value will retrieve the underlying interface value.
func (g *Goson) Value() interface{} {
	return g.i
}

// Len will return len() on the underlying value.
// If the value does not have a length the return value will be 0.
func (g *Goson) Len() int {
	switch t := g.i.(type) {
	case string:
		return len(t)
	case []interface{}:
		return len(t)
	case map[string]interface{}:
		return len(t)
	default:
		return 0
	}
}

// Index is used to access the index of an array object.
func (g *Goson) Index(index int) *Goson {
	if v, ok := g.i.([]interface{}); ok {
		return &Goson{i: v[index]}
	}
	return &Goson{i: new(interface{})}
}

// Bool returns the bool value.
func (g *Goson) Bool() bool {
	if v, ok := g.i.(bool); ok {
		return v
	}
	return false
}

// Int returns the underlying Int value converted from a float64.
func (g *Goson) Int() int {
	if v, ok := g.i.(float64); ok {
		return int(v)
	}
	return 0
}

// Float returns the underlying float64 value.
func (g *Goson) Float() float64 {
	if v, ok := g.i.(float64); ok {
		return v
	}
	return 0
}

// Slice returns the underlying slice value.
func (g *Goson) Slice() []interface{} {
	if v, ok := g.i.([]interface{}); ok {
		return v
	}
	return []interface{}{}
}

// Map returns the underlying map value.
func (g *Goson) Map() map[string]interface{} {
	if v, ok := g.i.(map[string]interface{}); ok {
		return v
	}
	return map[string]interface{}{}
}

// String will convert the underlying value as a string if it can or return an empty string.
// If the object is a JSON map or array it will return the structure in an indented format.
func (g *Goson) String() string {
	switch t := g.i.(type) {
	case bool:
		return strconv.FormatBool(t)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case string:
		return t
	case []interface{}:
		data, err := json.MarshalIndent(t, "", "\t")
		if err == nil {
			return string(data)
		}
	case map[string]interface{}:
		data, err := json.MarshalIndent(t, "", "\t")
		if err == nil {
			return string(data)
		}
	}
	return ""
}
