package service

import (
	"github.com/bwmarrin/snowflake"
	"github.com/pkg/errors"
)

const (
	// Epoch is set to the twitter snowflake epoch of Nov 04 2010 01:42:54 UTC in milliseconds
	// You may customize this to set a different epoch for your application.
	Epoch int64 = 1288834974657

	// NodeBits holds the number of bits to use for Node
	// Remember, you have a total 22 bits to share between Node/Step
	NodeBits uint8 = 10

	// StepBits holds the number of bits to use for Step
	// Remember, you have a total 22 bits to share between Node/Step
	StepBits  uint8 = 12
	nodeMax   int64 = -1 ^ (-1 << NodeBits)
	nodeMask        = nodeMax << StepBits
	stepMask  int64 = -1 ^ (-1 << StepBits)
	timeShift       = NodeBits + StepBits
	nodeShift       = StepBits
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

	if nodeNum < 0 || nodeNum > 1023 {
		return nil, errors.New("node out of range")
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

// Time returns an int64 unix timestamp in milliseconds of the snowflake ID time
func IDParseTime(id int64) int64 {
	return (id >> timeShift) + Epoch
}

// Node returns an int64 of the snowflake ID node number
func IDParseNode(id int64) int64 {
	return id & nodeMask >> nodeShift
}

// Step returns an int64 of the snowflake step (or sequence) number
func IDParseStep(id int64) int64 {
	return id & stepMask
}
