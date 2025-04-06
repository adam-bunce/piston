package packet

import (
	"fmt"
	"piston/internal"
	"reflect"
	"testing"
	"unsafe"
)

// TODO: test multiple with opaque boundaries in pipe
func TestDecode(t *testing.T) {
	cases := []struct {
		Name     string
		Input    []byte
		Expected map[string]any
	}{
		{
			Name:     "Handshake",
			Input:    []byte{2, 0, 5, 0, 95, 0, 52, 0, 100, 0, 97, 0, 109},
			Expected: map[string]any{"id": ID(Handshake), "username": []rune{'_', '4', 'd', 'a', 'm'}},
		},
		{
			Name:  "Login Request",
			Input: []byte{1, 0, 0, 0, 14, 0, 5, 0, 95, 0, 52, 0, 100, 0, 97, 0, 109},
			Expected: map[string]any{
				"id":              ID(LoginRequest),
				"protocolVersion": 14,
				"username":        []rune{'_', '4', 'd', 'a', 'm'},
				// assuming these will never be sent (maybe bad?, how would client even know this stuff right server is responsible for tracking that)
				// "mapSeed":         nil,
				// dimension: nil
			},
		},
		{
			Name:     "Player Position & Look",
			Input:    New(WithID(PlayerPositionAndLook), WithDouble(6.5), WithDouble(65.6), WithDouble(67.24), WithDouble(7.5), WithFloat(0.0), WithFloat(0.0), WithBool(false)).Body,
			Expected: map[string]any{"id": ID(PlayerPositionAndLook), "x": 6.5, "y": 65.6, "stance": 67.24, "z": 7.5, "yaw": float32(0.0), "pitch": float32(0.0), "onGround": false},
		},
	}

	for _, tc := range cases {
		fmt.Println()
		fmt.Println(tc.Name, "input size", unsafe.Sizeof(tc.Input))
		fmt.Println("->", tc.Input)
		client, server := internal.TestConn(t)

		client.Write(tc.Input)
		actual, err := ParsePacket(server)
		if err != nil {
			t.Fatalf("%s: %q", tc.Name, err)
		}

		client.Close()
		server.Close()

		if !reflect.DeepEqual(tc.Expected, actual) {
			for k, expectedVal := range tc.Expected {
				actualVal, ok := actual[k]
				if !ok {
					t.Logf("\t%q not present in actual", k)
				} else {
					if !reflect.DeepEqual(expectedVal, actualVal) {
						t.Logf("\tkey %q got %q (%T) expected %q (%T)", k, actualVal, actualVal, expectedVal, expectedVal)
					}
				}
			}

			t.Fatalf("%s: got %q expected %q", tc.Name, actual, tc.Expected)
		}
	}
}
