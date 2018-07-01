package BLC

import (
	"math/big"
	"crypto/sha256"
	"bytes"
)

// 256位Hash里面前面至少要有16个零(前导为4个0)
const targetBit = 20

type ProofOfWork struct {
	block *Block //当前要验证区块
	target *big.Int //大数据存储
}

//1.创建工作量对象
func NewProofOfWork(block *Block) *ProofOfWork {
	//1.big.Int对象 1
	// 2
	//0000 0001
	// 8 - 2 = 6
	// 0100 0000  64
	// 0010 0000
	// 0000 0000 0000 0001 0000 0000 0000 0000 0000 0000 .... 0000

	//1. 创建一个初始值为1的target
	target := big.NewInt(1)

	//2. 左移256 - targetBit
	target = target.Lsh(target,256 - targetBit)

	return &ProofOfWork{block,target}
}

//2.验证
func (proofOfWork *ProofOfWork) Run() ([]byte, int64){

	//1.将Block属性拼接成字节数组
	//2.生成hash
	//3.判断hash的有效性，如果满足条件则退出，不满足则一直循环求hash及对应的nonce值

	nonce := 0

	var hashInt big.Int //存储新生成的hash
	var hash [32]byte

	for {
		//准备数据
		dataBytes := proofOfWork.preparedData(nonce)
		//生成hash
		hash = sha256.Sum256(dataBytes)
		//fmt.Printf("\r%x",hash)

		//将hash存储到hashInt
		hashInt.SetBytes(hash[:])

		//判断hashInt是否小于Block里设置的target
		if 	proofOfWork.target.Cmp(&hashInt) == 1 {
			break
		}
		nonce = nonce + 1

	}

	return hash[:],int64(nonce)
}

//3.数据拼接，返回字节数组
func (proofOfWork *ProofOfWork) preparedData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			proofOfWork.block.PrevBlockHash,
			proofOfWork.block.Data,
			IntToHex(proofOfWork.block.Timestamp),
			IntToHex(int64(targetBit)),
			IntToHex(int64(nonce)),
			IntToHex(proofOfWork.block.Height),
		},
		[]byte{},
	)
	return data
}
