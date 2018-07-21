package BLC

import (
	"fmt"
	"net"
	"log"
	"bytes"
	"encoding/gob"
	"io"
	"io/ioutil"
	"encoding/hex"
)

const protocol = "tcp"   // 节点协议
const nodeVersion = 1    // 节点版本
const commandLength = 12 // 命令长度

var nodeAddress string                         // 当前节点地址
var miningAddress string                       // 挖矿地址
var knownNodes = []string{"localhost:3000"}    // 存储所有已知节点
var blocksInTransit = [][]byte{}               // 存储接收到的区块hash
var mempool = make(map[string]Transaction) // 存储接收到的交易

type yxh_addr struct {
	Yxh_AddrList []string
}

type yxh_block struct {
	Yxh_AddrFrom string
	Yxh_Block    []byte
}

type yxh_getblocks struct {
	Yxh_AddrFrom string
}

type yxh_getdata struct {
	Yxh_AddrFrom string
	Yxh_Type     string
	Yxh_ID       []byte
}

type yxh_inv struct {
	Yxh_AddrFrom string
	Yxh_Type     string
	Yxh_Items    [][]byte
}

type rwq_txs struct {
	Rwq_AddFrom     string
	Rwq_Transactions [][]byte
}

type yxh_version struct {
	Yxh_Version    int
	Yxh_BestHeight int
	Yxh_AddrFrom   string
}

//启动Server
func Yxh_StartServer(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	miningAddress = minerAddress
	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	bc := Yxh_NewBlockchain(nodeID)

	// 3000端口为：主节点
	// 3001端口为：钱包节点
	// 3002端口为：挖矿节点
	if nodeAddress != knownNodes[0] {
		// 此节点是钱包节点或者矿工节点，需要向主节点发送请求同步数据
		yxh_sendVersion(knownNodes[0], bc)
	}

	for { // 等待接收命令
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go yxh_handleConnecton(conn, bc)
	}
}

// 向中心节点发送 version 消息来查询是否自己的区块链已过时
func yxh_sendVersion(addr string, bc *Blockchain) {
	bestHeight := bc.Yxh_GetBestHeight()
	payload := yxh_gobEncode(yxh_version{nodeVersion, bestHeight, nodeAddress})

	request := append(yxh_commandToBytes("version"), payload...)

	yxh_sendData(addr, request)
}

// 发送数据
func yxh_sendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	// 如果连接失败，更新节点数据
	if err != nil {
		fmt.Sprintf("%s地址不可用\n", addr)
		var updatedNodes []string

		for _, node := range knownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}

		knownNodes = updatedNodes
		return
	}
	defer conn.Close()
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}

}

// 发送获取区块的的命令
func yxh_sendGetBlocks(address string) {
	payload := yxh_gobEncode(yxh_getblocks{nodeAddress})
	request := append(yxh_commandToBytes("getblocks"), payload...)

	yxh_sendData(address, request)
}

// 发送处理区块及交易hash列表请求
func rwq_sendInv(address, kind string, items [][]byte) {
	inventory := yxh_inv{nodeAddress, kind, items}
	payload := yxh_gobEncode(inventory)
	request := append(yxh_commandToBytes("inv"), payload...)

	yxh_sendData(address, request)
}

// 发送获取区块详细数据的命令
func yxh_sendGetData(address, kind string, id []byte) {
	payload := yxh_gobEncode(yxh_getdata{nodeAddress, kind, id})
	request := append(yxh_commandToBytes("getdata"), payload...)

	yxh_sendData(address, request)
}

// 发送区块内容命令
func yxh_sendBlock(addr string, b *Block) {
	data := yxh_block{nodeAddress, b.Yxh_Serialize()}
	payload := yxh_gobEncode(data)
	request := append(yxh_commandToBytes("block"), payload...)

	yxh_sendData(addr, request)
}

// 发送交易内容命令
func yxh_sendTx(addr string, tx *Transaction) {
	txs := []*Transaction{tx}
	yxh_sendTxs(addr,txs)
}
// 发送交易内容命令
func yxh_sendTxs(addr string, txs []*Transaction) {

	data := rwq_txs{nodeAddress, Yxh_SerializeTransactions(txs)}
	payload := yxh_gobEncode(data)
	request := append(yxh_commandToBytes("tx"), payload...)

	yxh_sendData(addr, request)
}

//================================================================
// 命令收集并分发
func yxh_handleConnecton(conn net.Conn, bc *Blockchain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := yxh_bytesToCommand(request[:commandLength])
	fmt.Printf("接收到命令：%s\n", command)

	switch command {
	case "addr": // 添加新节点
		yxh_handleAddr(request)
	case "block": // 添加新区块
		yxh_handleBlock(request, bc)
	case "inv": // 接收区块及交易hash列表 ，区块获取到内容到存储到 blocksInTransit ， 交易存储到 mempool
		yxh_handleInv(request, bc)
	case "getblocks": // 将当前节点区块链中的所有区块hash，返回给请求节点
		yxh_handleGetBlocks(request, bc)
	case "getdata": // 将单个交易或区块的内容 返回给请求节点
		yxh_handleGetData(request, bc)
	case "tx": // 添加新的交易,交易数量大于2，矿工节点挖矿,如果是主节点，进行分发交易
		yxh_handleTx(request, bc)
	case "version": // 检查是否需要同步数据，根据区块的height
		yxh_handleVersion(request, bc)
	default:
		fmt.Println("未知命令!")
	}

	conn.Close()

}

// 添加新节点
func yxh_handleAddr(request []byte) {
	var buff bytes.Buffer
	var payload yxh_addr

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	knownNodes = append(knownNodes, payload.Yxh_AddrList...)
	fmt.Printf("有%d个节点加入\n", len(knownNodes))
	// 把新节点推送给其他所有节点
	for _, node := range knownNodes {
		yxh_sendGetBlocks(node)
	}
}

/*
当接收到一个新块时，我们把它放到区块链里面。
如果还有更多的区块需要下载，我们继续从上一个下载的块的那个节点继续请求。
当最后把所有块都下载完后，对 UTXO 集进行重新索引

TODO: 并非无条件信任，我们应该在将每个块加入到区块链之前对它们进行验证。
TODO: 并非运行 UTXOSet.Reindex()， 而是应该使用 UTXOSet.Update(block)，
TODO: 因为如果区块链很大，它将需要很多时间来对整个 UTXO 集重新索引。
 */
func yxh_handleBlock(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload yxh_block

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Yxh_Block
	block := Yxh_DeserializeBlock(blockData)

	fmt.Println("收到一个新的区块!")
	bc.Yxh_AddBlock(block)

	fmt.Printf("Added block %x\n", block.Yxh_Hash)

	// 如果还有更多的区块需要下载，继续从上一个下载的块的那个节点继续请求
	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		yxh_sendGetData(payload.Yxh_AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := UTXOSet{bc}
		UTXOSet.Yxh_Reindex()
	}
}

// 向其他节点展示当前节点有什么块和交易,接收区块及交易列表
func yxh_handleInv(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload yxh_inv

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("接收到列表,长度为：%d，类型： %s\n", len(payload.Yxh_Items), payload.Yxh_Type)

	// 如果数据是 区块
	if payload.Yxh_Type == "block" {
		blocksInTransit = payload.Yxh_Items

		blockHash := payload.Yxh_Items[0]
		// 发送获取区块详细数据的命令
		yxh_sendGetData(payload.Yxh_AddrFrom, "block", blockHash)

		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}
	// 如果数据是 交易
	// 转账时，未立即挖矿，将交易保存到内存池中
	if payload.Yxh_Type == "tx" {
		txID := payload.Yxh_Items[0]
		// 如果内存池中，交易内容为空
		if mempool[hex.EncodeToString(txID)].Yxh_ID == nil {
			// 发送获取交易详细内容请求
			yxh_sendGetData(payload.Yxh_AddrFrom, "tx", txID)
		}
	}
}

// 处理获取区块命令，返回区块链中的所有区块hash
func yxh_handleGetBlocks(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload yxh_getblocks

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := bc.Yxh_GetBlockHashes()
	rwq_sendInv(payload.Yxh_AddrFrom, "block", blocks)
}

//  将单个交易或区块的内容 返回给请求节点
func yxh_handleGetData(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload yxh_getdata

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Yxh_Type == "block" {
		block, err := bc.Yxh_GetBlock([]byte(payload.Yxh_ID))
		if err != nil {
			return
		}

		yxh_sendBlock(payload.Yxh_AddrFrom, &block)
	}

	if payload.Yxh_Type == "tx" {
		txID := hex.EncodeToString(payload.Yxh_ID)
		tx := mempool[txID]

		yxh_sendTx(payload.Yxh_AddrFrom, &tx)
		// delete(mempool, txID)
	}
}

// 处理交易
func yxh_handleTx(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload rwq_txs

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	txData := payload.Rwq_Transactions
	txsDes := Yxh_DeserializeTransactions(txData)

	for _,tx := range txsDes {
		mempool[hex.EncodeToString(tx.Yxh_ID)] = tx

		// 如果是主节点
		if nodeAddress == knownNodes[0] {
			for _, node := range knownNodes {
				// 给其他节点分发，添加交易
				if node != nodeAddress && node != payload.Rwq_AddFrom {
					rwq_sendInv(node, "tx", [][]byte{tx.Yxh_ID})
				}
			}
		} else {
			// 如果交易池中有两条交易 并且当前是挖矿节点
			if len(mempool) >= 2 && len(miningAddress) > 0 {
			MineTransactions:
				var txs []*Transaction

				for id := range mempool {
					tx := mempool[id]
					if bc.Yxh_VerifyTransaction(&tx, txs) {
						txs = append(txs, &tx)
					}
				}

				if len(txs) == 0 {
					fmt.Println("交易不可用...")
					break
				}

				cbTx := Yxh_NewCoinbaseTX(miningAddress, "")
				txs = append(txs, cbTx)

				newBlock := bc.Yxh_MineBlock(txs)
				UTXOSet := UTXOSet{bc}
				UTXOSet.Update(newBlock)

				fmt.Println("挖到新的区块!")

				for _, tx := range txs {
					txID := hex.EncodeToString(tx.Yxh_ID)
					delete(mempool, txID)
				}

				for _, node := range knownNodes {
					if node != nodeAddress {
						rwq_sendInv(node, "block", [][]byte{newBlock.Yxh_Hash})
					}
				}

				if len(mempool) > 0 {
					goto MineTransactions
				}
			}
		}
	}
}

// 检查是否需要同步数据
func yxh_handleVersion(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload yxh_version
	// 获取数据
	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)

	if err != nil {
		log.Panic(err)
	}

	// 获取当前节点的最后height
	myBestHeight := bc.Yxh_GetBestHeight()
	// 数据中的最后height
	foreignerBestHeight := payload.Yxh_BestHeight

	// 节点将从消息中提取的 BestHeight 与自身进行比较
	// 当前的height比对方的小
	// 发送获取区块的的命令
	if myBestHeight < foreignerBestHeight {
		yxh_sendGetBlocks(payload.Yxh_AddrFrom)
	} else if myBestHeight > foreignerBestHeight {
		// 当前的height比对方的大
		// 通知对方节点，同步数据
		yxh_sendVersion(payload.Yxh_AddrFrom, bc)
	}

	// 如果不是已知节点，将节点添加到已知节点中
	if !yxh_nodeIsKnown(payload.Yxh_AddrFrom) {
		knownNodes = append(knownNodes, payload.Yxh_AddrFrom)
	}
}

// 是否是已知节点
func yxh_nodeIsKnown(addr string) bool {
	for _, node := range knownNodes {
		if node == addr {
			return true
		}
	}

	return false
}

//================================================================

// 命令转字节数组
func yxh_commandToBytes(command string) []byte {
	var bytes [commandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

// 将字节数组转字符串命令
func yxh_bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}

// 加密
func yxh_gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
