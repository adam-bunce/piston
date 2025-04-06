package packet

import (
	"reflect"
	"testing"
)

// might be good to sniff packets from actual server and use them as golden files
func TestStringToBytes(t *testing.T) {
	tt := []struct {
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

	for _, tc := range tt {
		actual := StringToBytes(tc.Input)

		if !reflect.DeepEqual(actual, tc.Expected) {
			t.Fatalf("%s: got '%v' expected '%v'", tc.Name, actual, tc.Expected)
		}
	}
}

func TestUtf16ToUtf8(t *testing.T) {
	tt := []struct {
		Name     string
		Input    []byte
		Expected []rune
	}{
		{
			Name:     "Username",
			Input:    []byte{0, 95, 0, 52, 0, 100, 0, 97, 0, 109},
			Expected: []rune{'_', '4', 'd', 'a', 'm'},
		},
	}

	for _, tc := range tt {
		actual := UTF16ToRunes(tc.Input)

		if !reflect.DeepEqual(actual, tc.Expected) {
			t.Fatalf("%s: got '%v' expected '%v'", tc.Name, actual, tc.Expected)
		}
	}

}

func TestLongToInt(t *testing.T) {
	tt := []struct {
		Name     string
		Input    []byte
		Expected int
	}{
		{
			Name:     "Username",
			Input:    []byte{0, 5, 0, 95, 0, 52, 0, 100},
			Expected: 1407782908854372,
		},
	}

	for _, tc := range tt {
		actual := LongToInt(tc.Input)

		if !reflect.DeepEqual(actual, tc.Expected) {
			t.Fatalf("%s: got '%v' expected '%v'", tc.Name, actual, tc.Expected)
		}
	}

}
