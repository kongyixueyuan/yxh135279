package BLC

import "bytes"

type TXInput struct {
	// 1. 交易的Hash
	Yxh_TxHash      []byte
	// 2. 存储TXOutput在Vout里面的索引
	Yxh_Vout      int

	Yxh_Signature []byte // 数字签名

	Yxh_PublicKey    []byte // 公钥，钱包里面
}

// 判断当前的消费是谁的钱
func (txInput *TXInput) Yxh_UnLockRipemd160Hash(ripemd160Hash []byte) bool {

	publicKey := Yxh_Ripemd160Hash(txInput.Yxh_PublicKey)

	return bytes.Compare(publicKey,ripemd160Hash) == 0
}