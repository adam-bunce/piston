package packet

import (
	"reflect"
	"testing"
)

func TestPacketBuilder(t *testing.T) {
	cases := []struct {
		Name      string
		Functions []func(*Packet)
		Expected  []byte
	}{
		{
			Name: "Test Login Request",
			Functions: []func(*Packet){
				WithID(LoginRequest),
				WithInt4(1),
				WithString16(""),
				WithLong(2),
				WithByte(0)},
			Expected: []byte{LoginRequest, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0},
		},
	}

	for _, tc := range cases {
		actual := New(tc.Functions...)

		if !reflect.DeepEqual(actual.Body, tc.Expected) {
			t.Fatalf("%s: got '%v' expected '%v'", tc.Name, actual, tc.Expected)
		}
	}

}
