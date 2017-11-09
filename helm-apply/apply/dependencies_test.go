package apply

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestResolveDependencies(t *testing.T) {
	cases := []struct {
		releases      map[string]*Release
		expectedOrder []string
		expectedErr   error
	}{
		{
			map[string]*Release{"1": {Name: "1"}},
			[]string{"1"},
			nil,
		},
		{
			map[string]*Release{
				"1": {Name: "1"},
				"2": {Name: "2", Dependencies: map[string]string{"dep": "1"}},
			},
			[]string{"1", "2"},
			nil,
		},
		{
			map[string]*Release{
				"1": {Name: "1", Dependencies: map[string]string{"dep": "2"}},
				"2": {Name: "2", Dependencies: map[string]string{"dep": "1"}},
			},
			nil,
			errors.New("dependency cycle detected"),
		},
	}
	for _, c := range cases {
		order, err := ResolveDependencies(c.releases)
		if !reflect.DeepEqual(order, c.expectedOrder) {
			t.Errorf("Expected order: %q but actual is: %q", c.expectedOrder, order)
		}
		if !(err == nil && c.expectedErr == nil) {
			if err == nil || c.expectedErr == nil || !strings.Contains(err.Error(), c.expectedErr.Error()) {
				t.Errorf("Expected err: %q but actual is: %q", c.expectedErr, err)
			}
		}
	}
}
