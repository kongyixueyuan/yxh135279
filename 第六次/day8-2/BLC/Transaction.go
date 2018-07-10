package BLC

import (
	"bytes"
	"log"
	"encoding/gob"
	"crypto/sha256"
	"encoding/hex"
	"crypto/ecdsa"
	"crypto/rand"

	"math/big"
	"crypto/elliptic"
	"time"
)

// UTXO
type Transaction struct {

	//1. 交易hash
	Yxh_TxHash []byte

	//2. 输入
	Yxh_Vins []*TXInput

	//3. 输出
	Yxh_Vouts []*TXOutput
}

//[]byte{}

// 判断当前的交易是否是Coinbase交易
func (tx *Transaction) Yxh_IsCoinbaseTransaction() bool {

	return len(tx.Yxh_Vins[0].Yxh_TxHash) == 0 && tx.Yxh_Vins[0].Yxh_Vout == -1
}

//1. Transaction 创建分两种情况
//1. 创世区块创建时的Transaction
func Yxh_NewCoinbaseTransaction(address string) *Transaction {

	//代表消费
	txInput := &TXInput{[]byte{},-1,nil,[]byte{}}

	txOutput := Yxh_NewTXOutput(10,address)

	txCoinbase := &Transaction{[]byte{},[]*TXInput{txInput},[]*TXOutput{txOutput}}

	//设置hash值
	txCoinbase.Yxh_HashTransaction()

	return txCoinbase
}

// 事务hash
func (tx *Transaction) Yxh_HashTransaction()  {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	resultBytes := bytes.Join([][]byte{Yxh_IntToHex(time.Now().Unix()),result.Bytes()},[]byte{})

	hash := sha256.Sum256(resultBytes)

	tx.Yxh_TxHash = hash[:]
}

//2. 转账时产生的Transaction
func Yxh_NewSimpleTransaction(from string,to string,amount int64,utxoSet *UTXOSet,txs []*Transaction) *Transaction {

	//$ ./bc send -from '["juncheng"]' -to '["zhangqiang"]' -amount '["2"]'
	//	[juncheng]
	//	[zhangqiang]
	//	[2]

	wallets,_ := Yxh_NewWallets()
	wallet := wallets.Yxh_WalletsMap[from]

	// 通过一个函数，返回
	money,spendableUTXODic := utxoSet.Yxh_FindSpendableUTXOS(from,amount,txs)
	//
	//	{hash1:[0],hash2:[2,3]}

	var txIntputs []*TXInput
	var txOutputs []*TXOutput

	for txHash,indexArray := range spendableUTXODic  {

		txHashBytes,_ := hex.DecodeString(txHash)
		for _,index := range indexArray  {

			txInput := &TXInput{txHashBytes,index,nil,wallet.Yxh_PublicKey}
			txIntputs = append(txIntputs,txInput)
		}

	}

	// 转账
	txOutput := Yxh_NewTXOutput(int64(amount),to)
	txOutputs = append(txOutputs,txOutput)

	// 找零
	txOutput = Yxh_NewTXOutput(int64(money) - int64(amount),from)
	txOutputs = append(txOutputs,txOutput)

	tx := &Transaction{[]byte{},txIntputs,txOutputs}

	//设置hash值
	tx.Yxh_HashTransaction()

	//进行签名
	utxoSet.Yxh_Blockchain.Yxh_SignTransaction(tx, wallet.Yxh_PrivateKey,txs)

	return tx

}

//产生Hash
func (tx *Transaction) Yxh_Hash() []byte {

	txCopy := tx

	txCopy.Yxh_TxHash = []byte{}

	hash := sha256.Sum256(txCopy.Yxh_Serialize())

	return hash[:]
}

//序列化
func (tx *Transaction) Yxh_Serialize() []byte {

	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)

	err := enc.Encode(tx)

	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

// 签名
func (tx *Transaction) Yxh_Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	//创世区块不用签名
	if tx.Yxh_IsCoinbaseTransaction() {
		return
	}

	for _, vin := range tx.Yxh_Vins {
		if prevTXs[hex.EncodeToString(vin.Yxh_TxHash)].Yxh_TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.Yxh_TrimmedCopy()

	for inID, vin := range txCopy.Yxh_Vins {
		prevTx := prevTXs[hex.EncodeToString(vin.Yxh_TxHash)]
		txCopy.Yxh_Vins[inID].Yxh_Signature = nil
		txCopy.Yxh_Vins[inID].Yxh_PublicKey = prevTx.Yxh_Vouts[vin.Yxh_Vout].Yxh_Ripemd160Hash
		txCopy.Yxh_TxHash = txCopy.Yxh_Hash()
		txCopy.Yxh_Vins[inID].Yxh_PublicKey = nil

		// 签名代码
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.Yxh_TxHash)
		if err != nil {
			log.Panic(err)
		}

		signature := append(r.Bytes(), s.Bytes()...)

		tx.Yxh_Vins[inID].Yxh_Signature = signature
	}
}

// 拷贝一份新的Transaction用于签名
func (tx *Transaction) Yxh_TrimmedCopy() Transaction {

	var inputs []*TXInput
	var outputs []*TXOutput

	for _, vin := range tx.Yxh_Vins {
		inputs = append(inputs, &TXInput{vin.Yxh_TxHash, vin.Yxh_Vout, nil, nil})
	}

	for _, vout := range tx.Yxh_Vouts {
		outputs = append(outputs, &TXOutput{vout.Yxh_Value, vout.Yxh_Ripemd160Hash})
	}

	txCopy := Transaction{tx.Yxh_TxHash, inputs, outputs}

	return txCopy
}

// 数字签名验证
func (tx *Transaction) Yxh_Verify(prevTXs map[string]Transaction) bool {

	if tx.Yxh_IsCoinbaseTransaction() {
		return true
	}

	for _, vin := range tx.Yxh_Vins {
		if prevTXs[hex.EncodeToString(vin.Yxh_TxHash)].Yxh_TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.Yxh_TrimmedCopy()

	curve := elliptic.P256()

	for inID, vin := range tx.Yxh_Vins {
		prevTx := prevTXs[hex.EncodeToString(vin.Yxh_TxHash)]
		txCopy.Yxh_Vins[inID].Yxh_Signature = nil
		txCopy.Yxh_Vins[inID].Yxh_PublicKey = prevTx.Yxh_Vouts[vin.Yxh_Vout].Yxh_Ripemd160Hash
		txCopy.Yxh_TxHash = txCopy.Yxh_Hash()
		txCopy.Yxh_Vins[inID].Yxh_PublicKey = nil

		// 私钥 ID
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Yxh_Signature)
		r.SetBytes(vin.Yxh_Signature[:(sigLen / 2)])
		s.SetBytes(vin.Yxh_Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.Yxh_PublicKey)
		x.SetBytes(vin.Yxh_PublicKey[:(keyLen / 2)])
		y.SetBytes(vin.Yxh_PublicKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.Yxh_TxHash, &r, &s) == false {
			return false
		}
	}

	return true
}
