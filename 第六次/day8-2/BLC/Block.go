package BLC

import (
	"time"
	"fmt"
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	//1. 区块高度
	Yxh_Height int64
	//2. 上一个区块HASH
	Yxh_PrevBlockHash []byte
	//3. 交易数据
	Yxh_Txs []*Transaction
	//4. 时间戳
	Yxh_Timestamp int64
	//5. Hash
	Yxh_Hash []byte
	// 6. Nonce
	Yxh_Nonce int64
}

// 需要将Txs转换成[]byte
func (block *Block) Yxh_HashTransactions() []byte  {
	//普通实现

	//var txHashes [][]byte
	//var txHash [32]byte
	//
	//for _, tx := range block.Txs {
	//	txHashes = append(txHashes, tx.TxHash)
	//}
	//txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	//
	//return txHash[:]

	//应用MerkleTree实现
	var transactions [][]byte

	for _, tx := range block.Yxh_Txs {
		transactions = append(transactions, tx.Yxh_Serialize())
	}
	mTree := Yxh_NewMerkleTree(transactions)

	return mTree.Yxh_RootNode.Yxh_Data

}

// 将区块序列化成字节数组
func (block *Block) Yxh_Serialize() []byte {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// 反序列化
func Yxh_DeserializeBlock(blockBytes []byte) *Block {

	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))

	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
	}

	return &block
}


//1. 创建新的区块
func Yxh_NewBlock(txs []*Transaction,height int64,prevBlockHash []byte) *Block {

	//创建区块
	block := &Block{height,prevBlockHash,txs,time.Now().Unix(),nil,0}

	// 调用工作量证明的方法并且返回有效的Hash和Nonce
	pow := Yxh_NewProofOfWork(block)

	// 挖矿验证
	hash,nonce := pow.Yxh_Run()

	block.Yxh_Hash = hash[:]

	block.Yxh_Nonce = nonce

	fmt.Println()

	return block
}

//2. 单独写一个方法，生成创世区块
func Yxh_CreateGenesisBlock(txs []*Transaction) *Block {

	return Yxh_NewBlock(txs,1, []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}

