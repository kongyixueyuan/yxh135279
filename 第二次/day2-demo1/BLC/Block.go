package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
	"fmt"
	"strconv"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

/*

网络上传输或能够保存至文件，它必须被编码然后再解码。当然，已经有许多可用的编码方式了：JSON，XML，Google 的 protocol buffers,又多了一种，由 Go 的 gob 包提供的方式
gob和json的pack之类的方法一样，由发送端使用Encoder对数据结构进行编码。在接收端收到消息之后，接收端使用Decoder将序列化的数据变化成本地变量
golang可以通过json或gob来序列化struct对象,虽然json的序列化更为通用,但利用gob编码可以实现json所不能支持的struct的方法序列化,利用gob包序列化struct保存到本地也十分简单
这里需要明确一点，gob只能用在golang中，所以在实际工程开发过程中，如果与其他端，或者其他语言打交道，那么gob是不可以的，我们就要使用json了

*/


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

//序列化 将Block转换成字节数组
func (block *Block) Serialize() []byte {
	//接收节点缓存
	var result bytes.Buffer
	//生成转换对象
	encoder := gob.NewEncoder(&result)
	//转换
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}
	//返回序列化后的字节数组
	return result.Bytes()
}


//反序列化，将字节数组转Block
func DeserializeBlock(blockBytes []byte) *Block {
	//结果接收对象
	var block Block
	//生成转换对象
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	//转换
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	//返回对象的地址，可以节点内存空间
	return &block

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