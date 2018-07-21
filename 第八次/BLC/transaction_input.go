package BLC

import "bytes"

type TXInput struct {
	Yxh_Txid      []byte
	Yxh_Vout      int      // Vout的index
	Yxh_Signature []byte   // 签名
	Yxh_PubKey    []byte   // 公钥
}

func (in TXInput) UsesKey(pubKeyHash []byte) bool  {
	lockingHash := Yxh_HashPubKey(in.Yxh_PubKey)

	return bytes.Compare(lockingHash,pubKeyHash) == 0
}
