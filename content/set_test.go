package content

import (
	"reflect"
	"sort"
	"testing"
)

func compareSlices(a, b []string) bool {
	sort.StringSlice(a).Sort()
	sort.StringSlice(b).Sort()
	return reflect.DeepEqual(a, b)
}

func TestStringSet(t *testing.T) {
	tests := []struct {
		Start []string
		Del   []string
		Add   []string
		End   []string
	}{
		{
			Start: []string{"one"},
			Del:   []string{"one"},
			Add:   []string{"two", "three"},
			End:   []string{"two", "three"},
		},
		{
			Start: []string{"one", "two"},
			Del:   []string{"one"},
			Add:   []string{"one", "two"},
			End:   []string{"one", "two"},
		},
		{
			Start: []string{""},
			Del:   []string{"three"},
			Add:   []string{"one", "two"},
			End:   []string{"one", "two"},
		},
		{
			Start: []string(nil),
			Del:   []string{"xxx"},
			Add:   []string{"one", "two"},
			End:   []string{"one", "two"},
		},
		{
			Start: []string(nil),
			Del:   []string(nil),
			Add:   []string(nil),
			End:   []string(nil),
		},
		{
			Start: []string{"one", "two"},
			Del:   []string{"two", "one"},
			Add:   []string(nil),
			End:   []string(nil),
		},
	}

	for i, tc := range tests {
		ss := NewStringSet(tc.Start)
		ss.RemoveAll(tc.Del)
		ss.AddAll(tc.Add)
		es := ss.Slice()
		if !compareSlices(tc.End, es) {
			t.Errorf("test %d failed: %#v != %#v", i, es, tc.End)
			t.Errorf("Test case:\n  %v\n- %v\n+ %v\n= %v", tc.Start, tc.Del, tc.Add, tc.End)
		}
	}

}
