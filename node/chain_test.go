package node

import (
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

func TestAddBlockWithTxInsufficientFunds(t *testing.T) {
	var (
		bs        = NewMemoryBlockStore()
		txs       = newMemoryTXStore()
		chain     = NewChain(bs, txs)
		block     = RandomBlock(t, chain)
		privKey   = crypto.NewPrivateKeyFromSeedStr(seed)
		recipient = crypto.GeneratePrivateKey().Public().Address().Bytes()
	)
	prevTx, err := chain.txStore.Get("8ada2924e739ee52ea194129ccc96ba93e9a87cbe465b912f381334cd7b939d0")

	assert.Nil(t, err)
	fmt.Println(prevTx)
	inputs := []*proto.TxInput{
		{
			PrevTxHash:   types.HashTransaction(prevTx),
			PrevOutIndex: 0,
			PublicKey:    privKey.Public().Bytes(),
		},
	}
	outputs := []*proto.TxOutput{
		{
			Amount:  10001,
			Address: recipient,
		},
	}
	tx := &proto.Transaction{
		Version: 1,
		Inputs:  inputs,
		Outputs: outputs,
	}

	sig := types.SignTransaction(privKey, tx)
	tx.Inputs[0].Signature = sig.Bytes()

	block.Transactions = append(block.Transactions, tx)
	require.NotNil(t, chain.AddBlock(block))
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
	prevTx, err := chain.txStore.Get("8ada2924e739ee52ea194129ccc96ba93e9a87cbe465b912f381334cd7b939d0")

	assert.Nil(t, err)
	fmt.Println(prevTx)
	inputs := []*proto.TxInput{
		{
			PrevTxHash:   types.HashTransaction(prevTx),
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

	sig := types.SignTransaction(privKey, tx)
	tx.Inputs[0].Signature = sig.Bytes()

	block.Transactions = append(block.Transactions, tx)

	require.Nil(t, chain.AddBlock(block))

}
