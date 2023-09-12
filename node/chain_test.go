package node

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/LDM-A/GoBlocker/crypto"
	"github.com/LDM-A/GoBlocker/proto"
	"github.com/LDM-A/GoBlocker/types"
	"github.com/LDM-A/GoBlocker/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func RandomBlock(t *testing.T, chain *Chain) *proto.Block {
	privKey := crypto.GeneratePrivateKey()
	block := util.RandomBlock()
	prevBlock, err := chain.GetBlockByHeight(chain.Height())
	require.Nil(t, err)
	block.Header.PreviousHash = types.HashBlock(prevBlock)
	types.SignBlock(privKey, block)
	return block
}

func TestNewChain(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), newMemoryTXStore())
	assert.Equal(t, 0, chain.Height())
	_, err := chain.GetBlockByHeight(0)

	assert.Nil(t, err)
}

func TestChainHeight(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), newMemoryTXStore())
	for i := 0; i < 100; i++ {
		b := RandomBlock(t, chain)

		require.Nil(t, chain.AddBlock(b))
		require.Equal(t, chain.Height(), i+1)
	}
}

func TestAddBlock(t *testing.T) {
	bs := NewMemoryBlockStore()
	txs := newMemoryTXStore()
	chain := NewChain(bs, txs)

	for i := 0; i < 100; i++ {
		block := RandomBlock(t, chain)
		blockHash := types.HashBlock(block)

		require.Nil(t, chain.AddBlock(block))
		fetchedBlock, err := chain.GetBlockByHash(blockHash)
		require.Nil(t, err)
		require.Equal(t, block, fetchedBlock)

		fetchedBlockByHeight, err := chain.GetBlockByHeight(i + 1)
		require.Nil(t, err)
		require.Equal(t, block, fetchedBlockByHeight)
	}

}

func TestAddBlockWithTx(t *testing.T) {
	var (
		bs        = NewMemoryBlockStore()
		txs       = newMemoryTXStore()
		chain     = NewChain(bs, txs)
		block     = RandomBlock(t, chain)
		privKey   = crypto.NewPrivateKeyFromSeedStr(seed)
		recipient = crypto.GeneratePrivateKey().Public().Address().Bytes()
	)
	ftt, err := chain.txStore.Get("b074e82904eaf4d97fb5cdddfcfd63b0930f72e16387df8d395331a7788a2936")

	assert.Nil(t, err)
	fmt.Println(ftt)
	inputs := []*proto.TxInput{
		{
			PrevTxHash:   types.HashTransaction(ftt),
			PrevOutIndex: 0,
			PublicKey:    privKey.Public().Bytes(),
		},
	}
	outputs := []*proto.TxOutput{
		{
			Amount:  100,
			Address: recipient,
		},
		{
			Amount:  900,
			Address: privKey.Public().Address().Bytes(),
		},
	}
	tx := &proto.Transaction{
		Version: 1,
		Inputs:  inputs,
		Outputs: outputs,
	}

	block.Transactions = append(block.Transactions, tx)

	require.Nil(t, chain.AddBlock(block))
	txHash := hex.EncodeToString(types.HashTransaction(tx))

	fetchedTx, err := chain.txStore.Get(txHash)
	assert.Nil(t, err)
	assert.Equal(t, tx, fetchedTx)

}
