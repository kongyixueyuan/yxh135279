package BLC

import "bytes"

type TXInput struct {
	// 1. 交易的Hash
	TxHash      []byte
	// 2. 存储TXOutput在Vout里面的索引
	Vout      int
	// 3. 数字签名
	Signature []byte
	// 4.用户公钥(原生)
	PubKey []byte
}

// 判断当前的消费是谁的钱
func (txInput *TXInput) UnLockRipemd160Hash(ripemd160Hash []byte) bool {

	publicKey := Ripemd160Hash(txInput.PubKey)

	return bytes.Compare(publicKey, ripemd160Hash) == 0
}