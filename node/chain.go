package node

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/LDM-A/GoBlocker/crypto"
	"github.com/LDM-A/GoBlocker/proto"
	"github.com/LDM-A/GoBlocker/types"
)

const seed = "f3c6d62c34725bd8c0c176738425d4d9e4a2f4d280886714f47e0acd250da504"

type HeaderList struct {
	headers []*proto.Header
}

func NewHeaderList() *HeaderList {
	return &HeaderList{
		headers: []*proto.Header{},
	}
}

func (list *HeaderList) Add(h *proto.Header) {
	list.headers = append(list.headers, h)
}

func (list *HeaderList) Len() int {
	return len(list.headers)
}

func (list *HeaderList) Height() int {
	return list.Len() - 1
}

type UTXO struct {
	Hash     string
	OutIndex int
	Amount   int64
	Spent    bool
}

type Chain struct {
	blockStore BlockStorer
	txStore    TXStorer
	utxoStore  UTXOStorer
	headers    *HeaderList
}

func NewChain(bs BlockStorer, txStore TXStorer) *Chain {
	chain := &Chain{
		blockStore: bs,
		txStore:    txStore,
		utxoStore:  newMemoryUTXOStore(),
		headers:    NewHeaderList(),
	}
	chain.addBlock(createGenesisBlock())

	return chain
}

func (list *HeaderList) Get(index int) *proto.Header {
	if index > list.Height() {
		panic("index too high")
	}
	return list.headers[index]
}

func (c *Chain) Height() int {
	return c.headers.Height()
}

func (c *Chain) AddBlock(b *proto.Block) error {
	if err := c.ValidateBlock(b); err != nil {
		return err
	}
	return c.addBlock(b)
}

func (c *Chain) addBlock(b *proto.Block) error {

	// Add the header to the list of headers
	c.headers.Add(b.Header)

	for _, tx := range b.Transactions {
		if err := c.txStore.Put(tx); err != nil {
			return err
		}
		hash := hex.EncodeToString(types.HashTransaction(tx))

		for it, output := range tx.Outputs {
			utxo := &UTXO{
				Hash:     hash,
				Amount:   output.Amount,
				OutIndex: it,
				Spent:    false,
			}
			if err := c.utxoStore.Put(utxo); err != nil {
				return err
			}
		}

	}

	return c.blockStore.Put(b)
}

func (c *Chain) GetBlockByHash(hash []byte) (*proto.Block, error) {
	hashHex := hex.EncodeToString(hash)
	return c.blockStore.Get(hashHex)
}

func (c *Chain) GetBlockByHeight(height int) (*proto.Block, error) {
	if c.Height() < height {
		return nil, fmt.Errorf("given height (%d) too high - height (%d)", height, c.Height())
	}
	header := c.headers.Get(height)
	hash := types.HashHeader(header)
	return c.GetBlockByHash(hash)
}

func (c *Chain) ValidateBlock(b *proto.Block) error {
	// validate sign of block
	if !types.VerifyBlock(b) {
		return fmt.Errorf("invalid block signature")
	}
	// Validate if the previous hash is actually the hash of the current block
	currentBlock, err := c.GetBlockByHeight(c.Height())
	if err != nil {
		return err
	}
	hash := types.HashBlock(currentBlock)
	if !bytes.Equal(hash, b.Header.PreviousHash) {
		return fmt.Errorf("invalid previous block hash")
	}

	for _, tx := range b.Transactions {
		if c.ValidateTransaction(tx); err != nil {
			return err
		}
	}
	return nil
}

func (c *Chain) ValidateTransaction(tx *proto.Transaction) error {
	if !types.VerifyTransaction(tx) {
		return fmt.Errorf("invalid transaction")
	}

	// check if all inputs are unspent by querying the utxo storage
	sumInputs := 0
	nInputs := len(tx.Inputs)
	hash := types.HashTransaction(tx)
	for i := 0; i < nInputs; i++ {
		key := fmt.Sprintf("%s_%d", tx.Inputs[i].PrevTxHash, i)
		utxo, err := c.utxoStore.Get(key)
		sumInputs += int(utxo.Amount)
		if err != nil {
			return err
		}
		if utxo.Spent {
			return fmt.Errorf("Output %d of %s is already spent", i, hash)
		}
	}
	sumOuts := 0
	for _, output := range tx.Outputs {
		sumOuts += int(output.Amount)
	}
	if sumInputs < sumOuts {
		return fmt.Errorf("Insufficient Balance")
	}
	return nil
}

func createGenesisBlock() *proto.Block {
	privKey := crypto.NewPrivateKeyFromSeedStr(seed)

	block := &proto.Block{
		Header: &proto.Header{
			Version: 1,
		},
	}

	tx := &proto.Transaction{
		Version: 1,
		Inputs:  []*proto.TxInput{},
		Outputs: []*proto.TxOutput{
			{
				Amount:  1000,
				Address: privKey.Public().Address().Bytes(),
			},
		},
	}

	block.Transactions = append(block.Transactions, tx)
	types.SignBlock(privKey, block)

	types.SignBlock(privKey, block)
	return block
}
