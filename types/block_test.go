package types

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/LDM-A/GoBlocker/proto"
	"github.com/LDM-A/GoBlocker/util"

	"github.com/LDM-A/GoBlocker/crypto"

	"github.com/stretchr/testify/assert"
)

func TestHashBlock(t *testing.T) {
	block := util.RandomBlock()
	hash := HashBlock(block)
	fmt.Println(hex.EncodeToString(hash))
	assert.Equal(t, 32, len(hash))
}

func TestVerifySignBlock(t *testing.T) {
	var (
		block   = util.RandomBlock()
		privKey = crypto.GeneratePrivateKey()
		pubkey  = privKey.Public()
	)

	sig := SignBlock(privKey, block)
	assert.Equal(t, 64, len(sig.Bytes()))
	assert.True(t, sig.Verify(pubkey, HashBlock(block)))

	assert.Equal(t, block.PublicKey, pubkey.Bytes())
	assert.Equal(t, block.Signature, sig.Bytes())

	assert.True(t, VerifyBlock(block))

	invalidPrivKey := crypto.GeneratePrivateKey()
	block.PublicKey = invalidPrivKey.Public().Bytes()

	assert.False(t, VerifyBlock(block))

}

func TestCalculateRootHash(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	block := util.RandomBlock()
	tx := &proto.Transaction{
		Version: 1,
	}
	block.Transactions = append(block.Transactions, tx)
	SignBlock(privKey, block)
	assert.True(t, VerifyRootHash(block))
	assert.Equal(t, 32, len(block.Header.RootHash))

}
