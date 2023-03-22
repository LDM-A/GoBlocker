package types

import (
	"fmt"
	"testing"

	"github.com/LarsDMsoftware/GoBlocker/crypto"
	"github.com/LarsDMsoftware/GoBlocker/proto"
	"github.com/LarsDMsoftware/GoBlocker/util"
)

func TestNewTransaction(t *testing.T) {

	fromPrivKey := crypto.GeneratePrivateKey()
	fromAddres := fromPrivKey.Public().Address().Bytes()

	toPrivKey := crypto.GeneratePrivateKey()
	toAddres := toPrivKey.Public().Address().Bytes()

	input := &proto.TxInput{
		PrevTxHash:   util.RandomHash(),
		PrevOutIndex: 0,
		PublicKey:    fromPrivKey.Public().Bytes(),
	}

	output1 := &proto.TxOutput{
		Amount:  5,
		Address: toAddres,
	}

	output2 := &proto.TxOutput{
		Amount:  95,
		Address: fromAddres,
	}

	tx := &proto.Transaction{
		Version: 1,
		Inputs:  []*proto.TxInput{input},
		Outputs: []*proto.TxOutput{output1, output2},
	}
	sig := SignTransaction(fromPrivKey, tx)
	input.Signature = sig.Bytes()

	fmt.Printf("%+v\n", tx)
}
