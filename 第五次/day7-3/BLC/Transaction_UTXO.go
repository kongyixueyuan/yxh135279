package BLC

type UTXO struct {
	TxHash []byte //交易hashID
	Index int  //在交易数据的位置
	Output *TXOutput //输出数据
}





