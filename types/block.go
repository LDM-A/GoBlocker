package types

import (
	"crypto/sha256"

	"github.com/LarsDMsoftware/GoBlocker/crypto"
	"github.com/LarsDMsoftware/GoBlocker/proto"

	pb "github.com/golang/protobuf/proto"
)

func SignBlock(pk *crypto.PrivateKey, b *proto.Block) *crypto.Signature {
	return pk.Sign(HashBlock(b))
}

// HashBlock returns a SHA256 of the header
func HashBlock(block *proto.Block) []byte {
	return HashHeader(block.Header)

}

func HashHeader(header *proto.Header) []byte {
	h, err := pb.Marshal(header)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(h)

	return hash[:]
}
