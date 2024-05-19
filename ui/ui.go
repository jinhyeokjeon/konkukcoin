package ui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"github.com/jinhyeokjeon/konkukcoin/blockchain"
	"github.com/jinhyeokjeon/konkukcoin/mine"
	"github.com/jinhyeokjeon/konkukcoin/p2p"
	"github.com/jinhyeokjeon/konkukcoin/utils"
	"github.com/jinhyeokjeon/konkukcoin/wallet"
	"golang.org/x/term"
)

func Start() {
	defer clearScreen()
	for {
		var choice string
		clearScreen()
		printKonkukCoin()
		fmt.Println()
		fmt.Print(" 1. 최신 블럭 조회     2. 모든 블럭 조회     3. 거래 조회     4.지갑 확인     5. 송금하기     6. 메모리 풀 조회\n\n")
		fmt.Print(" 7. 트랜잭션 검증      8. 검증한 트랜잭션 조회    9. 채굴     10. 노드 연결하기     11. 연결된 노드 확인     q. 종료\n\n")
		fmt.Printf(" >> ")
		fmt.Scanln(&choice)
		if choice == "q" {
			return
		}
		c, err := strconv.Atoi(choice)
		if err != nil {
			continue
		}
		switch c {
		case 1:
			latestBlock()
		case 2:
			allBlocks()
		case 3:
			findTx()
		case 4:
			checkWallet()
		case 5:
			transfer()
		case 6:
			checkMemPool()
		case 7:
			validateTx()
		case 8:
			validatedTxs()
		case 9:
			mining()
			fmt.Println()
			fmt.Printf(" 채굴이 완료되었습니다.\n")
		case 10:
			addNode()
		case 11:
			checkPeers()
		}
		proceed()
	}
}

func proceed() {
	fmt.Println()
	fmt.Print(" PRESS ENTER")
	oldState, _ := term.MakeRaw(int(syscall.Stdin))
	defer term.Restore(int(syscall.Stdin), oldState)
	for {
		var b [1]byte
		os.Stdin.Read(b[:])
		if b[0] == '\n' || b[0] == '\r' {
			break
		}
	}

}

func clearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	case "linux", "darwin": // "darwin"은 macOS를 의미합니다.
		cmd = exec.Command("clear")
	default:
		fmt.Println("Unsupported platform")
		return
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func printKonkukCoin() {
	fmt.Println()
	fmt.Println("        * * *                                                                                                             ")
	fmt.Println("     *         *                                                                                                          ")
	fmt.Println("   *   O     O   *    OOO       O     O     O     O     O     O     O     O       O O O       OOO       OOOOO     O     O ")
	fmt.Println("  *    O   O      * O     O     O O   O     O   O       O     O     O   O       O           O     O       O       O O   O ")
	fmt.Println("  *    O O        * O     O     O  O  O     O O         O     O     O O         O           O     O       O       O  O  O ")
	fmt.Println("  *    O   O      * O     O     O   O O     O   O       O     O     O   O       O           O     O       O       O   O O ")
	fmt.Println("   *   O     O   *    OOO       O    OO     O     O       O O       O     O       O O O       OOO       OOOOO     O    OO ")
	fmt.Println("     *         *                                                                                                          ")
	fmt.Println("        * * *                                                                                                             ")
}

func latestBlock() {
	lb := blockchain.Blockchain().LatestBlock
	if lb.Height == 0 {
		fmt.Print("\n 블록이 존재하지 않습니다.\n")
		return
	}
	fmt.Print(lb)
}

func allBlocks() {
	lb := blockchain.Blockchain().LatestBlock
	if lb.Height == 0 {
		fmt.Print("\n 블록이 존재하지 않습니다.\n")
		return
	}
	blocks := blockchain.Blockchain().GetAllBlocks()
	for i := len(blocks) - 1; i >= 0; i-- {
		fmt.Printf("%s", blocks[i])
	}
}

func findTx() {
	lb := blockchain.Blockchain().LatestBlock
	if lb.Height == 0 {
		fmt.Print("\n 블록이 존재하지 않습니다.\n")
		return
	}
	fmt.Printf("\n 조회할 거래 해시 입력: ")
	var hash string
	fmt.Scanln(&hash)
	fmt.Println()
	transactions := blockchain.Blockchain().GetAllTxs()

	tx := blockchain.Blockchain().GetTx(transactions, hash)
	if tx != nil {
		fmt.Print(tx.String(true))
		return
	}
	tx = blockchain.Blockchain().GetTxInMemPool(transactions, hash)
	if tx == nil {
		fmt.Printf("%s\n", " 거래가 존재하지 않습니다.")
	} else {
		fmt.Print(tx.String(false))
	}
}

func checkWallet() {
	w := wallet.Wallet()
	lb := blockchain.Blockchain().LatestBlock
	if lb.Height != 0 {
		unusedTxOuts := blockchain.Blockchain().UnusedTxOutsByAddress(w.Address)
		s := fmt.Sprintf(" %s\n", strings.Repeat("_", 152))
		s += utils.Print(152, 24, "address(public key)", wallet.Wallet().Address)
		s += utils.Print(152, 24, "", "")
		for index, tx := range unusedTxOuts {
			s += utils.Print(152, 24, fmt.Sprintf("[transaction_%d]", index), "")
			s += utils.Print(152, 24, "  tx_id", fmt.Sprint(tx.TxID))
			s += utils.Print(152, 24, "  output", fmt.Sprint(tx.OutputIdx))
			s += utils.Print(152, 24, "  amount", fmt.Sprint(tx.Amount))
			s += utils.Print(152, 24, "", "")
		}
		s += utils.Print(152, 24, "balance", fmt.Sprint(blockchain.Blockchain().GetBalance(unusedTxOuts)))
		s += fmt.Sprintf(" %s\n", strings.Repeat("-", 152))
		fmt.Print(s)
	} else {
		s := fmt.Sprintf(" %s\n", strings.Repeat("_", 152))
		s += utils.Print(152, 24, "address(public key)", wallet.Wallet().Address)
		s += utils.Print(152, 24, "balance", "0")
		s += fmt.Sprintf(" %s\n", strings.Repeat("-", 152))
		fmt.Print(s)
	}
}

func transfer() {
	lb := blockchain.Blockchain().LatestBlock
	if lb.Height == 0 {
		fmt.Print("\n 블록이 존재하지 않습니다.\n")
		return
	}
	addr := ""
	fmt.Print("\n 송금 받을 상대방의 주소를 입력하세요.\n\n")
	fmt.Print(" >> ")
	_, err := fmt.Scanln(&addr)
	if err != nil {
		utils.HandleErr(err)
	}
	balance := blockchain.Blockchain().GetBalance(blockchain.Blockchain().UnusedTxOutsByAddress(wallet.Wallet().Address))
	fmt.Print("\n 송금할 금액을 입력하세요. 현재 잔액은 ")
	fmt.Print(fmt.Sprint(balance) + " konkukcoin 입니다.\n\n")
	fmt.Print(" >> ")
	s := ""
	fmt.Scanln(&s)
	amount, err := strconv.Atoi(s)
	if err != nil {
		fmt.Print("\n 정수 형태로 올바르게 입력하여주세요.\n")
		return
	}
	if amount > balance {
		fmt.Print("\n 잔액이 부족합니다.\n")
		return
	}
	fmt.Print("\n 거래 수수료를 입력하세요. 송금액을 제외한 잔액은 ")
	fmt.Print(fmt.Sprint(balance-amount) + " konkukcoin 입니다.\n\n")
	fmt.Print(" >> ")
	fmt.Scanln(&s)
	fee, err := strconv.Atoi(s)
	if err != nil {
		fmt.Print("\n 정수 형태로 올바르게 입력하여주세요.\n")
		return
	}
	if fee+amount > balance {
		fmt.Print("\n 잔액이 부족합니다.\n")
		return
	}
	tx, total, _ := blockchain.AddTx(addr, amount, fee)
	p2p.BroadcastNewTx(tx)
	fmt.Println()
	s = fmt.Sprintf(" %s\n", strings.Repeat("_", 152))
	s += utils.Print(152, 24, "tx_id", tx.ID)
	s += utils.Print(152, 24, "from", wallet.Wallet().Address)
	s += utils.Print(152, 24, "to", addr)
	s += utils.Print(152, 24, "total used inputs", fmt.Sprintf("%s %s", fmt.Sprint(total), "konkuk coin"))
	s += utils.Print(152, 24, "transfer output", fmt.Sprintf("%s %s", fmt.Sprint(amount), "konkuk coin"))
	s += utils.Print(152, 24, "fee output", fmt.Sprintf("%s %s", fmt.Sprint(fee), "konkuk coin"))
	s += utils.Print(152, 24, "change output", fmt.Sprintf("%s %s", fmt.Sprint(total-amount-fee), "konkuk coin"))
	s += fmt.Sprintf(" %s\n", strings.Repeat("-", 152))
	fmt.Print(s)
}

func checkMemPool() {
	s := fmt.Sprintf(" %s\n", strings.Repeat("_", 88))
	index := 0
	for _, t := range blockchain.MemPool {
		if _, ok := blockchain.ConfirmedTxs[t.ID]; ok {
			continue
		}
		if index != 0 {
			s += utils.Print(88, 24, "", "")
		}
		s += utils.Print(88, 24, fmt.Sprintf("[tx_%d]", index), "")
		s += utils.Print(88, 24, "  tx_id", t.ID)
		s += utils.Print(88, 24, "  fee", fmt.Sprint(t.Fee))
		index++
	}
	if index == 0 {
		s += utils.Print(88, 40, "There is no Transaction.", "")
	}
	s += fmt.Sprintf(" %s\n", strings.Repeat("-", 88))
	fmt.Print(s)
}

func validateTx() {
	lb := blockchain.Blockchain().LatestBlock
	if lb.Height == 0 {
		fmt.Print("\n 블록이 존재하지 않습니다.\n")
		return
	}
	fmt.Print("\n 검증할 트랜잭션의 ID를 입력하세요.\n\n")
	fmt.Printf(" >> ")
	s := ""
	fmt.Scanln(&s)
	e := mine.AddTx(s)
	fmt.Println()
	switch e {
	case 0:
		fmt.Printf(" 트랜잭션은 최대 10개까지 검증 가능합니다.")
	case 1:
		fmt.Printf(" 메모리 풀 내에 존재하지 않는 트랜잭션입니다.")
	case 2:
		fmt.Printf(" 유효하지 않은 트랜잭션 입니다.")
	case 3:
		fmt.Printf(" 트랜잭션이 검증되었습니다.")
	}
	fmt.Println()
}

func validatedTxs() {
	s := fmt.Sprintf(" %s\n", strings.Repeat("_", 88))
	index := 0
	for _, t := range blockchain.ConfirmedTxs {
		if index != 0 {
			s += utils.Print(88, 24, "", "")
		}
		s += utils.Print(88, 24, fmt.Sprintf("[tx_%d]", index), "")
		s += utils.Print(88, 24, "  tx_id", t.ID)
		s += utils.Print(88, 24, "  fee", fmt.Sprint(t.Fee))
		index++
	}
	if index == 0 {
		s += utils.Print(88, 40, "There is no Transaction.", "")
	}
	s += fmt.Sprintf(" %s\n", strings.Repeat("-", 88))
	fmt.Print(s)
}

func mining() {
	oldState, _ := term.MakeRaw(int(syscall.Stdin))
	defer term.Restore(int(syscall.Stdin), oldState)
	fmt.Print("\033[?25l")
	newBlock := blockchain.Blockchain().AddBlock()
	fmt.Print("\033[?25h")
	p2p.BroadcastNewBlock(newBlock)
}

func addNode() {
	s := ""
	fmt.Println()
	fmt.Printf(" 연결할 노드 ip:port 입력\n\n")
	fmt.Print(" >> ")
	fmt.Scanln(&s)
	p2p.ConnectToPeer(s, true)
}

func checkPeers() {
	peers := p2p.AllPeers()
	if len(peers) == 0 {
		fmt.Printf("\n 연결된 노드가 없습니다.\n")
		return
	}
	for _, p := range peers {
		fmt.Printf("\n %s\n", p)
	}
}
