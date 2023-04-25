package node

import "github.com/LarsDMsoftware/GoBlocker/proto"

type BlockStorer interface {
	Put(*proto.Block) error
	Get(string)
}

type Chain struct {
	blockStore BlockStorer
}
