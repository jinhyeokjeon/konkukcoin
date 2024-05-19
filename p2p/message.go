package p2p

import (
	"bytes"
	"encoding/gob"

	"github.com/jinhyeokjeon/konkukcoin/blockchain"
	"github.com/jinhyeokjeon/konkukcoin/utils"
)

type MessageKind int

const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksRequest
	MessageAllBlocksResponse
	MessageNewBlockNotify
	MessageNewTxNotify
	MessageNewPeerNotify
)

type Message struct {
	Kind    MessageKind
	Payload []byte
}

func makeMessage(kind MessageKind, payload any) []byte {
	m := Message{
		Kind: kind,
	}
	if payload != nil {
		m.Payload = utils.ToBytes(payload)
	}
	return utils.ToBytes(m)
}

func sendNewestBlock(n *Node) {
	b := blockchain.Blockchain().LatestBlock
	m := makeMessage(MessageNewestBlock, b)
	n.inbox <- m
}

func requestAllBlocks(n *Node) {
	m := makeMessage(MessageAllBlocksRequest, nil)
	n.inbox <- m
}

func sendAllBlocks(n *Node) {
	if blockchain.Blockchain().LatestBlock.Height == 0 {
		return
	}
	m := makeMessage(MessageAllBlocksResponse, blockchain.Blockchain().GetAllBlocks())
	n.inbox <- m
}

func notifyNewBlock(b *blockchain.Block, n *Node) {
	m := makeMessage(MessageNewBlockNotify, b)
	n.inbox <- m
}

func notifyNewTx(tx *blockchain.Tx, n *Node) {
	m := makeMessage(MessageNewTxNotify, tx)
	n.inbox <- m
}

func notifyNewPeer(address string, n *Node) {
	m := makeMessage(MessageNewPeerNotify, address)
	n.inbox <- m
}

func handleMsg(m *Message, n *Node) {
	switch m.Kind {
	case MessageNewestBlock:
		var payload blockchain.Block
		decoder := gob.NewDecoder(bytes.NewReader(m.Payload))
		decoder.Decode(&payload)
		lb := blockchain.Blockchain().LatestBlock
		if lb.Height == 0 || payload.Height >= lb.Height {
			requestAllBlocks(n)
		} else {
			sendNewestBlock(n)
		}
	case MessageAllBlocksRequest:
		sendAllBlocks(n)
	case MessageAllBlocksResponse:
		var payload []*blockchain.Block
		decoder := gob.NewDecoder(bytes.NewReader(m.Payload))
		decoder.Decode(&payload)
		blockchain.Blockchain().Replace(payload)
	case MessageNewBlockNotify:
		var payload *blockchain.Block
		decoder := gob.NewDecoder(bytes.NewReader(m.Payload))
		decoder.Decode(&payload)
		blockchain.Blockchain().AddPeerBlock(payload)
	case MessageNewTxNotify:
		var payload *blockchain.Tx
		decoder := gob.NewDecoder(bytes.NewReader(m.Payload))
		decoder.Decode(&payload)
		blockchain.AddPeerTx(payload)
	case MessageNewPeerNotify:
		var payload string
		decoder := gob.NewDecoder(bytes.NewReader(m.Payload))
		decoder.Decode(&payload)
		ConnectToPeer(payload, false)
	}
}
