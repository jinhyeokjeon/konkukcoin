package blockchain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jinhyeokjeon/konkukcoin/utils"
	"github.com/jinhyeokjeon/konkukcoin/wallet"
)

const (
	minerReward int = 10
)

var MemPool = make(map[string]*Tx)
var ConfirmedTxs = make(map[string]Tx)

type Tx struct {
	ID        string
	Timestamp int64
	Inputs    []*TxIn
	Outputs   []*TxOut
	Fee       int
}

func (t *Tx) String(confirmed bool) string {
	s := fmt.Sprintf(" %s\n", strings.Repeat("_", 152))
	if confirmed {
		s += utils.Print(152, 24, "<CONFIRMED>", "")
	} else {
		s += utils.Print(152, 24, "<UNCONFIRMED>", "")
	}
	s += utils.Print(152, 24, "", "")
	s += utils.Print(152, 24, "tx_id", t.ID)
	s += utils.Print(152, 24, "time", utils.ConvertTime(t.Timestamp))
	for index, input := range t.Inputs {
		s += utils.Print(152, 24, "", "")
		s += utils.Print(152, 24, fmt.Sprintf("[input_%d]", index), "")
		s += input.String()
	}
	for index, output := range t.Outputs {
		s += utils.Print(152, 24, "", "")
		s += utils.Print(152, 24, fmt.Sprintf("[output_%d]", index), "")
		s += output.String()
	}
	s += fmt.Sprintf(" %s\n", strings.Repeat("-", 152))
	return s
}

type TxIn struct {
	TxID       string
	OutputIdx  int
	Signature  string
	KonkukBase bool
}

func (t *TxIn) String() string {
	s := ""
	s += utils.Print(152, 24, "  konkuk_base", fmt.Sprint(t.KonkukBase))
	if t.KonkukBase {
		s += utils.Print(152, 24, "  referenced_id", "nil")
		s += utils.Print(152, 24, "  output", "nil")
		s += utils.Print(152, 24, "  signature", "nil")
	} else {
		s += utils.Print(152, 24, "  referenced_id", t.TxID)
		s += utils.Print(152, 24, "  output", fmt.Sprint(t.OutputIdx))
		s += utils.Print(152, 24, "  signature", t.Signature)
	}
	return s
}

type TxOut struct {
	Address string
	Amount  int
}

func (t *TxOut) String() string {
	s := ""
	s += utils.Print(152, 24, "  address", t.Address)
	s += utils.Print(152, 24, "  amount", fmt.Sprint(t.Amount))
	return s
}

type UnusedTxOut struct {
	TxID      string
	OutputIdx int
	Amount    int
}

func (t *Tx) setId() {
	t.ID = utils.Hash(t)
}

func (t *Tx) sign() {
	for _, txIn := range t.Inputs {
		txIn.Signature = wallet.Sign(wallet.Wallet(), t.ID)
	}
}

func validate(tx *Tx) bool {
	txs := Blockchain().GetAllTxs()
	valid := true
	for _, txIn := range tx.Inputs {
		prevTx := Blockchain().GetTx(txs, txIn.TxID)
		if prevTx == nil {
			valid = false
			break
		}
		address := prevTx.Outputs[txIn.OutputIdx].Address
		valid = wallet.Verify(txIn.Signature, tx.ID, address)
		if !valid {
			break
		}
	}
	return valid
}

func isOnMempool(uTxOut *UnusedTxOut) bool {
	for _, tx := range MemPool {
		for _, input := range tx.Inputs {
			if input.TxID == uTxOut.TxID && input.OutputIdx == uTxOut.OutputIdx {
				return true
			}
		}
	}
	return false
}

func makeMinerRewardTx(address string) *Tx {
	txIns := []*TxIn{
		{"nil", -1, "KONKUK", true},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		ID:        "",
		Timestamp: time.Now().Unix(),
		Inputs:    txIns,
		Outputs:   txOuts,
		Fee:       0,
	}
	tx.setId()
	return &tx
}

func makeTx(from, to string, amount, fee int) (*Tx, int, error) {
	uTxOuts := Blockchain().UnusedTxOutsByAddress(from)

	if Blockchain().GetBalance(uTxOuts) < amount {
		return nil, -1, errors.New("not enough money")
	}

	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0

	for _, uTxOut := range uTxOuts {
		if total >= amount+fee {
			break
		}
		txIn := &TxIn{
			TxID:       uTxOut.TxID,
			OutputIdx:  uTxOut.OutputIdx,
			KonkukBase: false,
		}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}
	minerTxOut := &TxOut{"toMiner", fee}
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, minerTxOut)
	txOuts = append(txOuts, txOut)
	if change := total - amount - fee; change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}
	tx := &Tx{
		ID:        "",
		Timestamp: time.Now().Unix(),
		Inputs:    txIns,
		Outputs:   txOuts,
		Fee:       fee,
	}
	tx.setId()
	tx.sign()
	valid := validate(tx)
	if !valid {
		return nil, -1, errors.New("tx invalid")
	}
	return tx, total, nil
}

func AddTx(to string, amount, fee int) (*Tx, int, error) {
	tx, total, err := makeTx(wallet.Wallet().Address, to, amount, fee)
	if err != nil {
		return nil, -1, err
	}
	MemPool[tx.ID] = tx
	return tx, total, nil
}

func TxToConfirm() []*Tx {
	reward := makeMinerRewardTx(wallet.Wallet().Address)
	var txs []*Tx
	txs = append(txs, reward)
	for _, tx := range ConfirmedTxs {
		txs = append(txs, &tx)
	}
	return txs
}

func AddPeerTx(tx *Tx) {
	MemPool[tx.ID] = tx
}

func GetTxsInMemPool() []*Tx {
	var txs []*Tx
	for _, v := range MemPool {
		txs = append(txs, v)
	}
	return txs
}
