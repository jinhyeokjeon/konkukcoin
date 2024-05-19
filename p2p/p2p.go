package p2p

import (
	"github.com/jinhyeokjeon/konkukcoin/blockchain"
)

func BroadcastNewBlock(b *blockchain.Block) {
	for _, p := range peers {
		notifyNewBlock(b, p)
	}
}

func BroadcastNewTx(tx *blockchain.Tx) {
	for _, p := range peers {
		notifyNewTx(tx, p)
	}
}

func broadcastNewPeer(n *Node) {
	for key, p := range peers {
		if key != n.address {
			payload := n.address
			notifyNewPeer(payload, p)
		}
	}
}
