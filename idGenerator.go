package service

import (
	"github.com/bwmarrin/snowflake"
	"github.com/pkg/errors"
)

type idGenerator struct {
	gen *snowflake.Node
}

//node 0 - 1023
//+--------------------------------------------------------------------------+
//| 1 Bit Unused | 41 Bit Timestamp |  10 Bit NodeID  |   12 Bit Sequence ID |
//+--------------------------------------------------------------------------+
//Using the default settings, this allows for 4096 unique IDs
// to be generated every millisecond, per Node ID.
func newIDGenerator(nodeNum int64) (*idGenerator, error) {

	if nodeNum<0 || nodeNum>1023{
		return nil,errors.New("node out of range")
	}
	node, err := snowflake.NewNode(nodeNum)

	if err != nil {
		return nil, err
	}

	g := &idGenerator{
		gen: node,
	}

	return g, nil
}

func (g *idGenerator) GenerateInt64ID() int64 {

	return g.gen.Generate().Int64()
}
