package BLC

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
	"fmt"
)

type Block struct {
	Yxh_TimeStamp     int64
	Yxh_Transactions   []*Transaction
	Yxh_PrevBlockHash []byte
	Yxh_Hash          []byte
	Yxh_Nonce         int
	Yxh_Height        int
}
// 生成新的区块
func Yxh_NewBlock(transactions []*Transaction, prevBlockHash []byte, height int) *Block {
	// 生成新的区块对象
	block := &Block{
		time.Now().Unix(),
		transactions,
		prevBlockHash,
		[]byte{},
		0,
		height,
	}
	// 挖矿

	pow := Yxh_NewProofOfWork(block)
	nonce,hash :=pow.Yxh_Run()

	block.Yxh_Nonce = nonce
	block.Yxh_Hash = hash[:]

	return block

}

// 将交易进行hash
func (b Block) Yxh_HashTransactions() []byte {
	var transactions [][]byte
	// 获取交易真实内容
	for _,tx := range b.Yxh_Transactions{
		transactions = append(transactions,tx.Yxh_Serialize())
	}
	//txHash := sha256.Sum256(bytes.Join(transactions,[]byte{}))
	mTree := Yxh_NewMerkelTree(transactions)
	return mTree.Yxh_RootNode.Yxh_Data
}
// 新建创世区块
func Yxh_NewGenesisBlock(coinbase *Transaction) *Block  {
	return Yxh_NewBlock([]*Transaction{coinbase},[]byte{},1)
}

// 序列化
func (b *Block) Yxh_Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// 反序列化
func Yxh_DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
// 打印区块内容
func (block Block) String()  {
	fmt.Println("\n==============")
	fmt.Printf("Height:\t%d\n", block.Yxh_Height)
	fmt.Printf("PrevBlockHash:\t%x\n", block.Yxh_PrevBlockHash)
	fmt.Printf("Timestamp:\t%s\n", time.Unix(block.Yxh_TimeStamp, 0).Format("2006-01-02 03:04:05 PM"))
	fmt.Printf("Hash:\t%x\n", block.Yxh_Hash)
	fmt.Printf("Nonce:\t%d\n", block.Yxh_Nonce)
	fmt.Println("Txs:")

	for _, tx := range block.Yxh_Transactions {
		tx.String()
	}
	fmt.Println("==============")
}
