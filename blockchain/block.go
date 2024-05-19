package blockchain

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinhyeokjeon/konkukcoin/db"
	"github.com/jinhyeokjeon/konkukcoin/utils"
)

type Block struct {
	Hash       string
	PrevHash   string
	Height     int
	Difficulty int
	Nonce      int64
	Timestamp  int64
	Txs        []*Tx
}

func (b *Block) String() string {
	s := fmt.Sprintf(" %s\n", strings.Repeat("_", 88))
	s += utils.Print(88, 24, "hash", b.Hash)
	s += utils.Print(88, 24, "prev_hash", b.PrevHash)
	s += utils.Print(88, 24, "height", fmt.Sprint(b.Height))
	s += utils.Print(88, 24, "difficulty", fmt.Sprint(b.Difficulty))
	s += utils.Print(88, 24, "nonce", fmt.Sprint(b.Nonce))
	s += utils.Print(88, 24, "time", utils.ConvertTime(b.Timestamp))
	for index, t := range b.Txs {
		s += utils.Print(88, 24, "transaction_"+fmt.Sprint(index), t.ID)
	}
	s += fmt.Sprintf(" %s\n", strings.Repeat("-", 88))
	return s
}

func createBlock(prevHash string, height int, diff int) *Block {
	block := &Block{
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: diff,
		Nonce:      0,
	}
	block.Txs = TxToConfirm()
	block.mine()
	for _, tx := range ConfirmedTxs {
		delete(MemPool, tx.ID)
		delete(ConfirmedTxs, tx.ID)
	}
	saveBlock(block)
	return block
}

func (b *Block) mine() {
	s := ""
	target := strings.Repeat("0", b.Difficulty)
	fmt.Println()
	for {
		hash := utils.Hash(b)
		if s != "" {
			fmt.Print(strings.Repeat("\r", len(s)))
		}
		s = fmt.Sprintf(" target: %s / hash: %s / nonce: %s", strings.Repeat("0", b.Difficulty), hash, fmt.Sprint(b.Nonce))
		fmt.Print(s)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
	fmt.Println()
	b.Timestamp = time.Now().Unix()
}

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func saveBlock(b *Block) {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}
