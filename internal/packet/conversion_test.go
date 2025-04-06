package packet

import (
	"reflect"
	"testing"
)

// might be good to sniff packets from actual server and use them as golden files
func TestStringToBytes(t *testing.T) {
	cases := []struct {
		Name     string
		Input    string
		Expected []byte
	}{
		{
			Name:     "Username",
			Input:    "_4dam",
			Expected: []byte{0, 5, 0, 95, 0, 52, 0, 100, 0, 97, 0, 109},
		},
	}

	for _, tc := range cases {
		actual := StringToBytes(tc.Input)

		if !reflect.DeepEqual(actual, tc.Expected) {
			t.Fatalf("%s: got '%v' expected '%v'", tc.Name, actual, tc.Expected)
		}
	}
}
