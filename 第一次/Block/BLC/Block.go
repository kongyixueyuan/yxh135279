package BLC

import (
	"time"
	"fmt"
	"strconv"
	"bytes"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

type Block struct {
	//1.区块高度
	Height int64
	//2.上一次区块的HASH
	PrevBlockHash []byte
	//3.交易数据
	Data []byte
	//4.时间戳
	Timestamp int64
	//5.HASH
	Hash []byte
	//6.Nonce值
	Nonce int64
}


//1.创建新的区块
func NewBlock(data string, height int64, prevBlockHash []byte) *Block {
	//创建区块
	block := &Block{Height:height,PrevBlockHash:prevBlockHash,Data:[]byte(data),Timestamp:time.Now().Unix(),Hash:nil,Nonce:0}


	//设置挖矿数据
	pow := NewProofOfWork(block)
	//挖矿产生hash和nonce
	hash,nonce := pow.Run()

	//设置hash
	//block.SetHash()

	block.Hash = hash[:]
	block.Nonce = nonce

	fmt.Printf("	 hash:%x",block.Hash)
	fmt.Printf("	 nonce:%v",block.Nonce)

	fmt.Println()

	return block
}

//2.创建创世区块
func CreateGenesisBlock(data string) *Block {
	return NewBlock(data, 1, []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}


//3.设置Hash
func (block *Block) SetHash() {
	//生成最终的hash要求里面是字节数组，所以要将整形转换成字节数组
	//1. Height []byte
	heightBytes := IntToHex(block.Height)

	//2.时间戳转换成字节数组
	timeString := strconv.FormatInt(block.Timestamp,2)
	fmt.Println(timeString)
	timeBytes := []byte(timeString)

	//3.拼接所有属性(多个[]byte全成一个[]byte)
	blockBytes := bytes.Join([][]byte{heightBytes,block.PrevBlockHash,block.Data,timeBytes},[]byte{})

	//4.生成hash
	hash := sha3.Sum256(blockBytes)
	block.Hash = hash[:]

}