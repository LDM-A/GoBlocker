package types

import (
	"crypto/sha256"

	"github.com/LarsDMsoftware/GoBlocker/crypto"
	"github.com/LarsDMsoftware/GoBlocker/proto"

	pb "github.com/golang/protobuf/proto"
)

func SignTransaction(pk *crypto.PrivateKey, tx *proto.Transaction) *crypto.Signature {
	return pk.Sign(HashTransaction(tx))
}

func HashTransaction(tx *proto.Transaction) []byte {
	b, err := pb.Marshal(tx)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(b)

	return hash[:]
}
