package node

import (
	"testing"

	"github.com/LarsDMsoftware/GoBlocker/types"
	"github.com/LarsDMsoftware/GoBlocker/util"
	"github.com/stretchr/testify/assert"
)

func TestAddBlock(t *testing.T) {
	bs := NewMemoryBlockStore()
	chain := NewChain(bs)
	block := util.RandomBlock()
	blockHash := types.HashBlock(block)
	assert.Nil(t, chain.AddBlock(block))

	fetchedBlock, err := chain.GetBlockByHash(blockHash)
	assert.Nil(t, err)
	assert.Equal(t, block, fetchedBlock)
}
