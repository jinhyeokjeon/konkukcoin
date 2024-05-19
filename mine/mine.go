package mine

import (
	"github.com/jinhyeokjeon/konkukcoin/blockchain"
	"github.com/jinhyeokjeon/konkukcoin/wallet"
)

func validate(txId string) (int, *blockchain.Tx) {
	var address string
	var amount int = 0
	txs := blockchain.Blockchain().GetAllTxs()
	tx := blockchain.Blockchain().GetTx(blockchain.GetTxsInMemPool(), txId)
	if tx == nil {
		return 1, nil
	}
	valid := true
	for _, txIn := range tx.Inputs {
		prevTx := blockchain.Blockchain().GetTx(txs, txIn.TxID)
		if prevTx == nil {
			valid = false
			break
		}
		address = prevTx.Outputs[txIn.OutputIdx].Address
		valid = wallet.Verify(txIn.Signature, tx.ID, address)
		if !valid {
			break
		}
		amount += prevTx.Outputs[txIn.OutputIdx].Amount
	}
	if !valid {
		return 2, nil
	}
	tx.Outputs[0].Address = wallet.Wallet().Address
	return 3, tx
}

func AddTx(txId string) int {
	if len(blockchain.ConfirmedTxs) == 10 {
		return 0 // 최대 트랜잭션 개수 도달
	}
	e, tx := validate(txId)
	if e != 3 {
		return e // 1: 존재하지 않는 트랜잭션, 2: Invalid Transaction
	}
	blockchain.ConfirmedTxs[txId] = *tx
	return 3 // 성공적으로 추가.
}
