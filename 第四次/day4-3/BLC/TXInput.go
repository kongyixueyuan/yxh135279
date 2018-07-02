package BLC

type TXInput struct {
	//交易hash
	TxHash []byte
	//此笔交易在原交易区块中的位置
	Vout int
	//用户名, 某个人的钱
	ScriptPubKey string
}


// 判断当前的消费是谁的钱
func (txInput *TXInput) UnLockWithAddress(address string) bool {

	return txInput.ScriptPubKey == address
}