package p2p

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"
)

var Port string

type Node struct {
	address string
	conn    net.Conn
	inbox   chan []byte
}

var peers = make(map[string]*Node)

func Start() {
	ln, err := net.Listen("tcp", "localhost:"+Port)
	if err != nil {
		fmt.Println("Error starting node:", err)
		os.Exit(1)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		b := make([]byte, 10)
		conn.Read(b)
		addr := conn.RemoteAddr().String()
		idx := strings.Index(addr, ":")
		handleConnection(conn, addr[:idx], string(b))
	}
}

func handleConnection(conn net.Conn, ip, port string) *Node {
	address := ip + ":" + port
	fmt.Println("\n Connected to peer:", address)
	node := &Node{address, conn, make(chan []byte)}
	peers[address] = node
	go node.read()
	go node.write()
	return node
}

func (n *Node) close() {
	defer n.conn.Close()
	delete(peers, n.address)
}

func (n *Node) read() {
	defer n.close()
	defer close(n.inbox) // I added it
	for {
		m := Message{}
		decoder := gob.NewDecoder(n.conn)
		err := decoder.Decode(&m)
		if err != nil {
			break
		}
		handleMsg(&m, n)
	}
}

func (n *Node) write() {
	for {
		m, ok := <-n.inbox
		if !ok {
			break
		}
		n.conn.Write(m)
	}
}

func ConnectToPeer(address string, broadcast bool) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to peer:", err)
		return
	}
	conn.Write([]byte(Port))
	idx := strings.Index(address, ":")
	node := handleConnection(conn, address[:idx], address[idx+1:])
	if broadcast {
		broadcastNewPeer(node)
	}
	sendNewestBlock(node)
}

func AllPeers() []string {
	var keys []string
	for key := range peers {
		keys = append(keys, key)
	}
	return keys
}
