package BLC

import "bytes"

//TXOutput{100,"zhangbozhi"}
//TXOutput{30,"xietingfeng"}
//TXOutput{40,"zhangbozhi"}


type TXOutput struct {
	Value int64
	Ripemd160Hash []byte  //经过ripemd160后的公钥 与Input中的PubKey 经过计算后的结果对比是一致的
}

// 解锁
func (txOutput *TXOutput) UnLockScriptPubKeyWithAddress(address string) bool {

	publicKeyHash := Base58Decode([]byte(address))

	ripemd160Hash := publicKeyHash[1:len(publicKeyHash)-4]

	return bytes.Compare(txOutput.Ripemd160Hash, ripemd160Hash) == 0
}

// 通过地址得到公钥
func (txOutput *TXOutput) Lock(address string) {

	publicKeyHash := Base58Decode([]byte(address))
	//取中间的20个长度
	txOutput.Ripemd160Hash = publicKeyHash[1:len(publicKeyHash)-4]
}

// 产生新的输出对象
func NewTXOutput(value int64, address string) *TXOutput {

	txOutput := &TXOutput{value,nil}

	txOutput.Lock(address)

	return txOutput
}


