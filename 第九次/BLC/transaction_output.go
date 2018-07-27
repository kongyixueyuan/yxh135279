package BLC

import "bytes"

type TXOutput struct {
	Yxh_Value  int
	Yxh_PubKeyHash []byte
}
// 根据地址获取 PubKeyHash
func (out *TXOutput) Yxh_Lock(address []byte) {
	pubKeyHash := Yxh_Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.Yxh_PubKeyHash = pubKeyHash
}

// 判断是否当前公钥对应的交易输出(是否是某个人的交易输出)
func (out *TXOutput) Yxh_IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.Yxh_PubKeyHash, pubKeyHash) == 0
}

func Yxh_NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Yxh_Lock([]byte(address))
	return txo
}


