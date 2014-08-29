package packet

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var confs = []map[string]interface{}{
	{
		"pack":  Packet{Code: Req, Type: 3, Arguments: [][]byte{{4}, {5}, {6}}},
		"bytes": []byte{0, 0x52, 0x45, 0x51, 0, 0, 0, 3, 0, 0, 0, 5, 4, 0, 5, 0, 6},
	},
	{
		"pack":  Packet{Code: Req, Type: 3, Arguments: [][]byte{}},
		"bytes": []byte{0, 0x52, 0x45, 0x51, 0, 0, 0, 3, 0, 0, 0, 0},
	},
	{
		"pack":  Packet{Code: Req, Type: 3, Arguments: [][]byte{{4}}},
		"bytes": []byte{0, 0x52, 0x45, 0x51, 0, 0, 0, 3, 0, 0, 0, 1, 4},
	},
}

func TestBytes(t *testing.T) {
	for _, conf := range confs {
		pack := conf["pack"].(Packet)
		bytes, err := pack.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, bytes, conf["bytes"])
	}
}

func TestConstructor(t *testing.T) {
	for _, conf := range confs {
		pack := &Packet{}
		if err := pack.UnmarshalBinary(conf["bytes"].([]byte)); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, *pack, conf["pack"])
	}
}
