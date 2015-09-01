package packet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var confs = []struct {
	Pack  Packet
	Bytes []byte
}{
	{
		Pack:  Packet{Code: Req, Type: 3, Arguments: [][]byte{{4}, {5}, {6}}},
		Bytes: []byte{0, 0x52, 0x45, 0x51, 0, 0, 0, 3, 0, 0, 0, 5, 4, 0, 5, 0, 6},
	},
	{
		Pack:  Packet{Code: Req, Type: 3, Arguments: [][]byte{}},
		Bytes: []byte{0, 0x52, 0x45, 0x51, 0, 0, 0, 3, 0, 0, 0, 0},
	},
	{
		Pack:  Packet{Code: Req, Type: 3, Arguments: [][]byte{{4}}},
		Bytes: []byte{0, 0x52, 0x45, 0x51, 0, 0, 0, 3, 0, 0, 0, 1, 4},
	},
	{
		/*
			00 52 45 53                \0RES        (Magic)
			00 00 00 0b                11           (Packet type: JOB_ASSIGN)
			00 00 00 14                20           (Packet length)
			48 3a 6c 61 70 3a 31 00    H:lap:1\0    (Job handle)
			72 65 76 65 72 73 65 00    reverse\0    (Function)
			74 65 73 74                test         (Workload)
		*/
		Pack: Packet{
			Code: Res,
			Type: JobAssign,
			Arguments: [][]byte{
				[]byte("H:lap:1"),
				[]byte("reverse"),
				[]byte("test"),
			},
		},
		Bytes: []byte{
			0x00, 0x52, 0x45, 0x53,
			0x00, 0x00, 0x00, 0x0b,
			0x00, 0x00, 0x00, 0x14,
			0x48, 0x3a, 0x6c, 0x61, 0x70, 0x3a, 0x31, 0x00,
			0x72, 0x65, 0x76, 0x65, 0x72, 0x73, 0x65, 0x00,
			0x74, 0x65, 0x73, 0x74,
		},
	},
	{
		/*
			00 52 45 51                \0REQ        (Magic)
			00 00 00 07                7            (Packet type: SUBMIT_JOB)
			00 00 00 0d                13           (Packet length)
			72 65 76 65 72 73 65 00    reverse\0    (Function)
			00                         \0           (Unique ID)
			74 65 73 74                test         (Workload)
		*/
		Pack: Packet{
			Code: Req,
			Type: SubmitJob,
			Arguments: [][]byte{
				[]byte("reverse"),
				[]byte{},
				[]byte("test"),
			},
		},
		Bytes: []byte{
			0x00, 0x52, 0x45, 0x51,
			0x00, 0x00, 0x00, 0x07,
			0x00, 0x00, 0x00, 0x0d,
			0x72, 0x65, 0x76, 0x65, 0x72, 0x73, 0x65, 0x00,
			0x00,
			0x74, 0x65, 0x73, 0x74,
		},
	},
	{
		/*
			00 52 45 51                \0REQ        (Magic)
			00 00 00 04                4            (Packet type: PRE_SLEEP)
			00 00 00 00                0            (Packet length)
		*/
		Pack: Packet{
			Code:      Req,
			Type:      PreSleep,
			Arguments: [][]byte{},
		},
		Bytes: []byte{
			0x00, 0x52, 0x45, 0x51,
			0x00, 0x00, 0x00, 0x04,
			0x00, 0x00, 0x00, 0x00,
		},
	},
}

func TestBytes(t *testing.T) {
	for _, conf := range confs {
		bytes, err := conf.Pack.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, bytes, conf.Bytes)
	}
}

func TestConstructor(t *testing.T) {
	for _, conf := range confs {
		pack := &Packet{}
		if err := pack.UnmarshalBinary(conf.Bytes); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, *pack, conf.Pack)
	}
}
