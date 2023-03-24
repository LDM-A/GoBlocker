package node

import (
	"context"
	"fmt"
	"net"

	"github.com/LarsDMsoftware/GoBlocker/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type Node struct {
	peers map[net.Addr]*grpc.ClientConn
	proto.UnimplementedNodeServer
	version string
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error) {
	peer, _ := peer.FromContext(ctx)
	fmt.Println("Received tx from:", peer)
	return &proto.Ack{}, nil
}

func (n *Node) Handshake(ctx context.Context, v *proto.Version) (*proto.Version, error) {
	ourVersion := &proto.Version{
		Version: n.version,
		Height:  100,
	}
	return ourVersion, nil
}

func NewNode() *Node {
	return &Node{
		version: "v0.1",
	}
}
