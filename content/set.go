package content

import (
	"strings"
)

type StringSet struct {
	Values map[string]string
}

func NewStringSet(slice []string) StringSet {
	set := StringSet{Values: make(map[string]string)}
	if slice != nil {
		set.AddAll(slice)
	}
	return set
}

func (s StringSet) AddAll(slice []string) {
	if slice == nil {
		return
	}
	for _, v := range slice {
		k := strings.ToLower(v)
		s.Values[k] = v
	}
}

func (s StringSet) RemoveAll(slice []string) {
	if slice == nil {
		return
	}
	for _, v := range slice {
		k := strings.ToLower(v)
		delete(s.Values, k)
	}
}

func (s StringSet) Slice() []string {
	slice := make([]string, 0, len(s.Values))
	for _, v := range s.Values {
		if v != "" {
			slice = append(slice, v)
		}
	}
	if len(slice) < 1 || (len(slice) == 1 && slice[0] == "") {
		return nil
	}
	return slice
}
