package BLC

import (
	"math/big"
	"math"
	"bytes"
	"crypto/sha256"
	"fmt"
)

var (
	maxNonce = math.MaxInt64
)

const targetBits = 16

type ProofOfWork struct {
	Yxh_block  *Block
	Yxh_target *big.Int
}

// 生成新的工作量证明
func Yxh_NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}
	return pow
}

// 准备挖矿hash数据
func (pow *ProofOfWork) Yxh_PrepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.Yxh_block.Yxh_PrevBlockHash,
		pow.Yxh_block.Yxh_HashTransactions(),
		IntToHex(pow.Yxh_block.Yxh_TimeStamp),
		IntToHex(int64(targetBits)),
		IntToHex(int64(nonce)),
	}, []byte{})
	return data
}

// 执行工作量证明，返回nonce值和hash
func (pow *ProofOfWork) Yxh_Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte

	nonce := 0
	for nonce < maxNonce {
		data := pow.Yxh_PrepareData(nonce)

		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		//if math.Remainder(float64(nonce),100000) == 0{
		//	fmt.Printf("\r%x",hash)
		//}
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.Yxh_target) == -1 {
			break;
		} else {
			nonce++
		}
	}
	return nonce, hash[:]

}
