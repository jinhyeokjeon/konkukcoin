package blockchain

import (
	"errors"

	"github.com/jinhyeokjeon/konkukcoin/db"
	"github.com/jinhyeokjeon/konkukcoin/utils"
)

const (
	defaultDifficulty  int = 3
	difficultyInterval int = 3
	miningTimePerBlock int = 1
	allowedRange       int = 2
)

type blockchain struct {
	LatestBlock *Block
}

var b = &blockchain{&Block{Height: 0}}

func Blockchain() *blockchain {
	if b.LatestBlock.Height == 0 {
		hash := db.GetNewestHash()
		if hash != nil {
			b.restore(hash)
		}
	}
	return b
}

func (b *blockchain) restore(hash []byte) {
	var newestHash string
	utils.FromBytes(&newestHash, hash)
	b.LatestBlock, _ = b.GetBlock(newestHash)
}

func (b *blockchain) AddBlock() *Block {
	if b.LatestBlock == nil {
		b.LatestBlock = &Block{Hash: "", Height: 0}
	}
	block := createBlock(b.LatestBlock.Hash, b.LatestBlock.Height+1, b.getDifficulty())
	*b.LatestBlock = *block
	db.SaveNewestHash(utils.ToBytes(block.Hash))
	return block
}

var ErrNotFound = errors.New("block not found")

func (b *blockchain) GetBlock(hash string) (*Block, error) {
	blockBytes := db.GetBlock(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}

func (b *blockchain) GetAllBlocks() []*Block {
	var blocks []*Block
	hashCursor := b.LatestBlock.Hash
	for {
		block, _ := b.GetBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

func (b *blockchain) getDifficulty() int {
	if b.LatestBlock.Height <= difficultyInterval {
		return defaultDifficulty
	} else if b.LatestBlock.Height%difficultyInterval == 0 {
		return b.recalculateDifficulty()
	} else {
		return b.LatestBlock.Difficulty
	}
}

func (b *blockchain) recalculateDifficulty() int {
	allBlocks := b.GetAllBlocks()
	prevBlock := allBlocks[difficultyInterval]
	actualTime := (int)((b.LatestBlock.Timestamp / 60) - (prevBlock.Timestamp / 60))
	expectedTime := difficultyInterval * miningTimePerBlock
	if actualTime < expectedTime-allowedRange {
		return b.LatestBlock.Difficulty + 1
	} else if actualTime > expectedTime+allowedRange {
		return b.LatestBlock.Difficulty - 1
	}
	return b.LatestBlock.Difficulty
}

func (b *blockchain) GetAllTxs() []*Tx {
	var txs []*Tx
	for _, block := range b.GetAllBlocks() {
		txs = append(txs, block.Txs...)
	}
	return txs
}

func (b *blockchain) GetTx(transactions []*Tx, ID string) *Tx {
	for _, tx := range transactions {
		if tx.ID == ID {
			return tx
		}
	}
	return nil
}

func (b *blockchain) GetTxInMemPool(transactions []*Tx, ID string) *Tx {
	for _, tx := range MemPool {
		if tx.ID == ID {
			return tx
		}
	}
	return nil
}

func (b *blockchain) UnusedTxOutsByAddress(address string) []*UnusedTxOut {
	var uTxOuts []*UnusedTxOut
	transactions := b.GetAllTxs()
	usedTxs := make(map[string]bool)

	for _, block := range b.GetAllBlocks() {
		for _, tx := range block.Txs {
			for _, input := range tx.Inputs {
				if input.Signature == "KONKUK" {
					break
				}
				if b.GetTx(transactions, input.TxID).Outputs[input.OutputIdx].Address == address {
					usedTxs[input.TxID] = true
				}
			}
			if _, ok := usedTxs[tx.ID]; !ok {
				for index, output := range tx.Outputs {
					if output.Address == address {
						utxOut := &UnusedTxOut{tx.ID, index, output.Amount}
						if !isOnMempool(utxOut) {
							uTxOuts = append(uTxOuts, utxOut)
						}
					}
				}
			}
		}
	}
	return uTxOuts
}

func (b *blockchain) GetBalance(uTxOuts []*UnusedTxOut) int {
	amount := 0
	for _, txOut := range uTxOuts {
		amount += txOut.Amount
	}
	return amount
}

func (b *blockchain) Replace(newBlocks []*Block) {
	b.LatestBlock = newBlocks[0]
	db.SaveNewestHash(utils.ToBytes(b.LatestBlock.Hash))
	db.EmptyBlocks()
	for _, block := range newBlocks {
		saveBlock(block)
	}
}

func (b *blockchain) AddPeerBlock(newBlock *Block) {
	b.LatestBlock = newBlock
	db.SaveNewestHash(utils.ToBytes(b.LatestBlock.Hash))
	saveBlock(newBlock)

	for _, tx := range newBlock.Txs {
		_, ok := MemPool[tx.ID]
		if ok {
			delete(MemPool, tx.ID)
		}
		_, ok = ConfirmedTxs[tx.ID]
		if ok {
			delete(ConfirmedTxs, tx.ID)
		}
	}
}
