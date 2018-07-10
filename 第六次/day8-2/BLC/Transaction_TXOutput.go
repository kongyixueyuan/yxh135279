package BLC

import "bytes"

type TXOutput struct {
	Yxh_Value int64
	Yxh_Ripemd160Hash []byte  //用户名
}

func (txOutput *TXOutput)  Yxh_Lock(address string)  {

	publicKeyHash := Yxh_Base58Decode([]byte(address))

	txOutput.Yxh_Ripemd160Hash = publicKeyHash[1:len(publicKeyHash) - 4]
}


func Yxh_NewTXOutput(value int64,address string) *TXOutput {

	txOutput := &TXOutput{value,nil}

	// 设置Ripemd160Hash
	txOutput.Yxh_Lock(address)

	return txOutput
}

// 解锁
func (txOutput *TXOutput) Yxh_UnLockScriptPubKeyWithAddress(address string) bool {

	publicKeyHash := Yxh_Base58Decode([]byte(address))

	hash160 := publicKeyHash[1:len(publicKeyHash) - 4]

	return bytes.Compare(txOutput.Yxh_Ripemd160Hash,hash160) == 0
}



