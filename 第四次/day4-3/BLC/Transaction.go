package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
)

type Transaction struct {
	//每笔交易都有hash
	//注意，交易的Hash是通过所有输入和输出就能算出的，不是通过挖矿产生的
	TxHash []byte

	//交易输入 &TXInput{TxHash, pos, address}
	//用来记录转账过程中某一个人的某张钱被消费了
	Vins []*TXInput

	//交易输出 &TXOutput{100,address}
	Vouts []*TXOutput

}

func (tx *Transaction) HashTransaction() {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(buff.Bytes())
	tx.TxHash = hash[:]
}

//Coinbase帐号的Transaction
func NewCoinbaseTransaction(address string) *Transaction {
	//记录已经消费
	txInput := &TXInput{TxHash:[]byte{}, Vout:-1,ScriptPubKey:""}
	//记录已消费和未消费的支出
	txOutput := &TXOutput{Value:100,ScriptPubKey:address}
	//组织一条交易记录
	txCoinbase := &Transaction{TxHash:[]byte{}, Vins:[]*TXInput{txInput}, Vouts:[]*TXOutput{txOutput}}

	//设置hash
	txCoinbase.HashTransaction()

	return txCoinbase

}

//组织交易数据
//2. 转账时产生的Transaction

func NewSimpleTransaction(from string,to string,amount int,blockchain *Blockchain,txs []*Transaction) *Transaction {

	//$ ./bc send -from '["juncheng"]' -to '["zhangqiang"]' -amount '["2"]'
	//	[juncheng]
	//	[zhangqiang]
	//	[2]

	// 通过一个函数，返回
	money,spendableUTXODic := blockchain.FindSpendableUTXOS(from,amount,txs)
	//
	//	{hash1:[0],hash2:[2,3]}

	var txIntputs []*TXInput
	var txOutputs []*TXOutput

	for txHash,indexArray := range spendableUTXODic  {

		txHashBytes,_ := hex.DecodeString(txHash)
		for _,index := range indexArray  {
			txInput := &TXInput{txHashBytes,index,from}
			txIntputs = append(txIntputs,txInput)
		}

	}

	// 转账
	txOutput := &TXOutput{int64(amount),to}
	txOutputs = append(txOutputs,txOutput)

	// 找零
	txOutput = &TXOutput{int64(money) - int64(amount),from}
	txOutputs = append(txOutputs,txOutput)

	tx := &Transaction{[]byte{},txIntputs,txOutputs}

	//设置hash值
	tx.HashTransaction()

	return tx

}

// 判断当前的交易是否是Coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool {

	return len(tx.Vins[0].TxHash) == 0 && tx.Vins[0].Vout == -1
}
