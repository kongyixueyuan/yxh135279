package BLC

import (
	"bytes"
	"log"
	"encoding/gob"
	"crypto/sha256"
	"encoding/hex"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/elliptic"
	"math/big"
)

// UTXO
type Transaction struct {

	//1. 交易hash
	TxHash []byte

	//2. 输入
	Vins []*TXInput

	//3. 输出
	Vouts []*TXOutput
}

//[]byte{}

// 判断当前的交易是否是Coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool {

	return len(tx.Vins[0].TxHash) == 0 && tx.Vins[0].Vout == -1
}



//1. Transaction 创建分两种情况
//1. 创世区块创建时的Transaction
func NewCoinbaseTransaction(address string) *Transaction {

	//代表消费输入
	txInput := &TXInput{[]byte{},-1,nil, []byte{}}

	//输出
	txOutput := NewTXOutput(10, address)

	//创世区块交易数据
	txCoinbase := &Transaction{[]byte{},[]*TXInput{txInput},[]*TXOutput{txOutput}}

	//设置hash值
	txCoinbase.HashTransaction()

	return txCoinbase
}


// 计算交易hash
func (tx *Transaction) HashTransaction()  {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)

	if err != nil {
		log.Panic(err)
	}

	hash := sha256.Sum256(result.Bytes())

	tx.TxHash = hash[:]
}


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

	wallets, _ := NewWallets()
	wallet := wallets.Wallets

	for txHash,indexArray := range spendableUTXODic  {

		txHashBytes,_ := hex.DecodeString(txHash)

		for _,index := range indexArray  {
			txInput := &TXInput{txHashBytes,index, nil,wallet[from].PublicKey}
			txIntputs = append(txIntputs,txInput)
		}

	}

	// 转账
	txOutput := NewTXOutput(int64(amount),to)
	txOutputs = append(txOutputs,txOutput)

	// 找零
	txOutput = NewTXOutput(int64(money) - int64(amount),from)
	txOutputs = append(txOutputs,txOutput)

	tx := &Transaction{[]byte{},txIntputs,txOutputs}

	//设置hash值
	tx.HashTransaction()

	//签名
	blockchain.SignTransaction(tx, wallet[from].PrivateKey)

	return tx

}

//签名方法
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {

	//创世区块不用判断
	if tx.IsCoinbaseTransaction() {
		return
	}

	for _,vin := range tx.Vins {
		if prevTXs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panic("Previous tandsction is not correct")
		}
	}

	//复制交易数据，除去部分数据
	txCopy := tx.TrimmedCopy()

	for inID, Vin := range txCopy.Vins {
		prevTX := prevTXs[hex.EncodeToString(Vin.TxHash)]
		txCopy.Vins[inID].Signature = nil
		txCopy.Vins[inID].PubKey = prevTX.Vouts[Vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inID].PubKey = nil

		//签名
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.TxHash)

		if err != nil {
			log.Panic(err)
		}

		//组成签名的数据
		signature := append(r.Bytes(), s.Bytes() ...)

		tx.Vins[inID].Signature = signature
	}

}

// 部分数据复制
func (tx *Transaction) TrimmedCopy() Transaction {

	var inputs []*TXInput
	var outputs []*TXOutput

	for _, vin := range tx.Vins {
		//注意将签名和公钥除去，后面采用从交易区块中查询交易数据，得到output数据里的pubkey
		inputs = append(inputs, &TXInput{vin.TxHash, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vouts{
		//输出数据不变化
		outputs = append(outputs, &TXOutput{vout.Value, vout.Ripemd160Hash})
	}

	txCopy := Transaction{tx.TxHash, inputs, outputs}

	return txCopy

}

// 生成交易的hash
func (tx *Transaction) Hash() []byte {

	txCopy := tx

	txCopy.TxHash = []byte{}

	hash := sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

// 交易对象序列化
func (tx *Transaction) Serialize() []byte {

	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

//数字签名验证校验
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {

	//创世区块不用验证
	if tx.IsCoinbaseTransaction() {
		return true
	}

	for _, vin := range tx.Vins {

		if prevTXs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panic("error : previous transction is nog correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	curve := elliptic.P256()

	for inID, vin := range tx.Vins {

		prevTX := prevTXs[hex.EncodeToString(vin.TxHash)]
		txCopy.Vins[inID].Signature = nil
		txCopy.Vins[inID].PubKey = prevTX.Vouts[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inID].PubKey = nil

		//私钥
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:sigLen/2])
		s.SetBytes(vin.Signature[sigLen/2:])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:keyLen/2])
		y.SetBytes(vin.PubKey[keyLen/2:])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}

		if ecdsa.Verify(&rawPubKey, txCopy.TxHash, &r, &s) == false {
			return false
		}

	}

	return true

}