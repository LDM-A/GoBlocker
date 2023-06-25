package main

import (
	"context"
	"log"
	"time"

	"github.com/LDM-A/GoBlocker/crypto"
	"github.com/LDM-A/GoBlocker/node"
	"github.com/LDM-A/GoBlocker/proto"
	"github.com/LDM-A/GoBlocker/util"
	"google.golang.org/grpc"
)

func main() {
	// Create nodes 3000, 4000, 5000
	// Only connect 3000 with 4000 and 4000 with 5000
	// Through peer discovery 3000 will find and connect to 5000
	makeNode(":3000", []string{})
	time.Sleep(time.Second)
	makeNode(":4000", []string{":3000"})
	time.Sleep(time.Second)
	makeNode(":5000", []string{":4000"})

	time.Sleep(time.Second)

	makeTransaction()
	select {}
}

func makeNode(listenAddr string, bootstrapNodes []string) *node.Node {
	n := node.NewNode()
	go n.Start(listenAddr, bootstrapNodes)

	return n
}

func makeTransaction() {
	client, err := grpc.Dial(":3000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	c := proto.NewNodeClient(client)

	privKey := crypto.GeneratePrivateKey()
	tx := &proto.Transaction{
		Version: 1,
		Inputs: []*proto.TxInput{
			{
				PrevTxHash:   util.RandomHash(),
				PrevOutIndex: 0,
				PublicKey:    privKey.Public().Bytes(),
			},
		},
		Outputs: []*proto.TxOutput{
			{
				Amount:  99,
				Address: privKey.Public().Address().Bytes(),
			},
		},
	}

	_, err = c.HandleTransaction(context.TODO(), tx)
	if err != nil {
		log.Fatal(err)
	}
}
