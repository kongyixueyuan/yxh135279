package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"fmt"
	"strings"
	"encoding/hex"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/elliptic"
	"math/big"
)

// 创世区块，Token数量
const subsidy  = 10

type Transaction struct {
	Yxh_ID   []byte
	Yxh_Vin  []TXInput
	Yxh_Vout []TXOutput
}

// 是否是创世区块交易
func (tx Transaction) Yxh_IsCoinbase() bool {
	// Vin 只有一条
	// Vin 第一条数据的Txid 为 0
	// Vin 第一条数据的Vout 为 -1
	return len(tx.Yxh_Vin) == 1 && len(tx.Yxh_Vin[0].Yxh_Txid) == 0 && tx.Yxh_Vin[0].Yxh_Vout == -1
}


// 将交易进行Hash
func (tx *Transaction) Yxh_Hash() []byte  {
	var hash [32]byte

	txCopy := *tx
	txCopy.Yxh_ID = []byte{}

	hash = sha256.Sum256(txCopy.Yxh_Serialize())
	return hash[:]
}
// 新建创世区块的交易
func Yxh_NewCoinbaseTX(to ,data string) *Transaction  {
	if data == ""{
		//如果数据为空，可以随机给默认数据,用于挖矿奖励
		randData := make([]byte, 20)
		_, err := rand.Read(randData)
		if err != nil {
			log.Panic(err)
		}

		data = fmt.Sprintf("%x", randData)
	}
	txin := TXInput{[]byte{},-1,nil,[]byte(data)}
	txout := Yxh_NewTXOutput(subsidy,to)

	tx := Transaction{nil,[]TXInput{txin},[]TXOutput{*txout}}
	tx.Yxh_ID = tx.Yxh_Hash()
	return &tx
}

// 转帐时生成交易
func Yxh_NewUTXOTransaction(wallet *Wallet,to string,amount int,UTXOSet *UTXOSet,txs []*Transaction) *Transaction   {

	// 如果本区块中，多笔转账
	/**
	第一种情况：
	  A:10
	  A->B 2
	  A->C 4

	  tx1:
	      Vin:
	           ATxID  out ...
	      Vout:
	           A : 8
	           B : 2
	  tx1:
	      Vin:
	           ATxID  out ...
	      Vout:
	           A : 4
	           C : 4
	第二种情况：
	  A:10+10
	  A->B 4
	  A->C 8
	**/

	pubKeyHash := Yxh_HashPubKey(wallet.Yxh_PublicKey)
	if len(txs) > 0 {
		// 查的txs中的UTXO
		utxo := Yxh_FindUTXOFromTransactions(txs)

		// 找出当前钱包已经花费的
		unspentOutputs := make(map[string][]int)
		acc := 0
		for txID,outs := range utxo {
			for outIdx, out := range outs.Yxh_Outputs {
				if out.Yxh_IsLockedWithKey(pubKeyHash) && acc < amount {
					acc += out.Yxh_Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}

		if acc >= amount { // 当前交易中的剩余余额可以支付
			fmt.Println("txs>0 && acc >= amount")
			return Yxh_NewUTXOTransactionEnd(wallet,to,amount,UTXOSet,acc,unspentOutputs,txs)
		}else{
			fmt.Println("txs>0 && acc < amount")
			accLeft, validOutputs := UTXOSet.Yxh_FindSpendableOutputs(pubKeyHash,  amount - acc)
			for k,v := range unspentOutputs{
				validOutputs[k] = v
			}
			return Yxh_NewUTXOTransactionEnd(wallet,to,amount,UTXOSet,acc + accLeft,validOutputs,txs)
		}
	} else { //只是当前一笔交易
		fmt.Println("txs==0")
		acc, validOutputs := UTXOSet.Yxh_FindSpendableOutputs(pubKeyHash, amount)

		return Yxh_NewUTXOTransactionEnd(wallet,to,amount,UTXOSet,acc,validOutputs,txs)
	}
}

func Yxh_NewUTXOTransactionEnd(wallet *Wallet,to string,amount int,UTXOSet *UTXOSet,acc int,UTXO map[string][]int,txs []*Transaction) *Transaction {

	if acc < amount {
		log.Panic("账户余额不足")
	}

	var inputs []TXInput
	var outputs []TXOutput
	// 构造input
	for txid, outs := range UTXO {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{txID, out, nil, wallet.Yxh_PublicKey}
			inputs = append(inputs, input)
		}
	}
	// 生成交易输出
	outputs = append(outputs, *Yxh_NewTXOutput(amount, to))
	// 生成余额
	if acc > amount {
		outputs = append(outputs, *Yxh_NewTXOutput(acc-amount, string(wallet.Yxh_GetAddress())))
	}

	tx := Transaction{nil, inputs, outputs}
	tx.Yxh_ID = tx.Yxh_Hash()
	// 签名

	//tx.String()
	UTXOSet.Yxh_Blockchain.Yxh_SignTransaction(&tx, wallet.Yxh_PrivateKey,txs)

	return &tx
}


// 找出交易中的utxo
func Yxh_FindUTXOFromTransactions(txs []*Transaction) map[string]TXOutputs {
	UTXO := make(map[string]TXOutputs)
	// 已经花费的交易txID : TXOutputs.index
	spentTXOs := make(map[string][]int)
	// 循环区块中的交易
	for _, tx := range txs {
		// 将区块中的交易hash，转为字符串
		txID := hex.EncodeToString(tx.Yxh_ID)

	Outputs:
		for outIdx, out := range tx.Yxh_Vout { // 循环交易中的 TXOutputs
			// Was the output spent?
			// 如果已经花费的交易输出中，有此输出，证明已经花费
			if spentTXOs[txID] != nil {
				for _, spentOutIdx := range spentTXOs[txID] {
					if spentOutIdx == outIdx { // 如果花费的正好是此笔输出
						continue Outputs // 继续下一次循环
					}
				}
			}

			outs := UTXO[txID] // 获取UTXO指定txID对应的TXOutputs
			outs.Yxh_Outputs = append(outs.Yxh_Outputs, out)
			UTXO[txID] = outs
		}

		if tx.Yxh_IsCoinbase() == false { // 非创世区块
			for _, in := range tx.Yxh_Vin {
				inTxID := hex.EncodeToString(in.Yxh_Txid)
				spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Yxh_Vout)
			}
		}
	}
	return UTXO

}

// 签名
func (tx *Transaction) Yxh_Sign(privateKey ecdsa.PrivateKey,prevTXs map[string]Transaction)  {
	if tx.Yxh_IsCoinbase() { // 创世区块不需要签名
		return
	}

	// 检查交易的输入是否正确
	for _,vin := range tx.Yxh_Vin{
		if prevTXs[hex.EncodeToString(vin.Yxh_Txid)].Yxh_ID == nil{
			log.Panic("错误：之前的交易不正确")
		}
	}

	txCopy := tx.Yxh_TrimmedCopy()

	for inID, vin := range txCopy.Yxh_Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Yxh_Txid)]
		txCopy.Yxh_Vin[inID].Yxh_Signature = nil
		txCopy.Yxh_Vin[inID].Yxh_PubKey = prevTx.Yxh_Vout[vin.Yxh_Vout].Yxh_PubKeyHash

		dataToSign := fmt.Sprintf("%x\n", txCopy)

		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, []byte(dataToSign))
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Yxh_Vin[inID].Yxh_Signature = signature
		txCopy.Yxh_Vin[inID].Yxh_PubKey = nil
	}
}
// 验证签名
func (tx *Transaction) Yxh_Verify(prevTXs map[string]Transaction) bool {
	if tx.Yxh_IsCoinbase() {
		return true
	}

	for _, vin := range tx.Yxh_Vin {
		if prevTXs[hex.EncodeToString(vin.Yxh_Txid)].Yxh_ID == nil {
			log.Panic("错误：之前的交易不正确")
		}
	}

	txCopy := tx.Yxh_TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.Yxh_Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Yxh_Txid)]
		txCopy.Yxh_Vin[inID].Yxh_Signature = nil
		txCopy.Yxh_Vin[inID].Yxh_PubKey = prevTx.Yxh_Vout[vin.Yxh_Vout].Yxh_PubKeyHash

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Yxh_Signature)
		r.SetBytes(vin.Yxh_Signature[:(sigLen / 2)])
		s.SetBytes(vin.Yxh_Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.Yxh_PubKey)
		x.SetBytes(vin.Yxh_PubKey[:(keyLen / 2)])
		y.SetBytes(vin.Yxh_PubKey[(keyLen / 2):])

		dataToVerify := fmt.Sprintf("%x\n", txCopy)

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) == false {
			return false
		}
		txCopy.Yxh_Vin[inID].Yxh_PubKey = nil
	}

	return true
}

// 复制交易（输入的签名和公钥置为空）
func (tx *Transaction) Yxh_TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Yxh_Vin {
		inputs = append(inputs, TXInput{vin.Yxh_Txid, vin.Yxh_Vout, nil, nil})
	}

	for _, vout := range tx.Yxh_Vout {
		outputs = append(outputs, TXOutput{vout.Yxh_Value, vout.Yxh_PubKeyHash})
	}

	txCopy := Transaction{tx.Yxh_ID, inputs, outputs}

	return txCopy
}
// 打印交易内容
func (tx Transaction) String()  {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction ID: %x", tx.Yxh_ID))

	for i, input := range tx.Yxh_Vin {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Yxh_Txid))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Yxh_Vout))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Yxh_Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.Yxh_PubKey))
	}

	for i, output := range tx.Yxh_Vout {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Yxh_Value))
		lines = append(lines, fmt.Sprintf("       PubKeyHash: %x", output.Yxh_PubKeyHash))
	}
	fmt.Println(strings.Join(lines, "\n"))
}


// 将交易序列化
func (tx Transaction) Yxh_Serialize() []byte  {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)

	if err != nil{
		log.Panic(err)
	}
	return encoded.Bytes()
}
// 反序列化交易
func Yxh_DeserializeTransaction(data []byte) Transaction {
	var transaction Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}

	return transaction
}

// 将交易数组序列化
func Yxh_SerializeTransactions(txs []*Transaction) [][]byte  {

	var txsHash [][]byte
	for _,tx := range txs{
		txsHash = append(txsHash, tx.Yxh_Serialize())
	}
	return txsHash
}

// 反序列化交易数组
func Yxh_DeserializeTransactions(data [][]byte) []Transaction {
	var txs []Transaction
	for _,tx := range data {
		var transaction Transaction
		decoder := gob.NewDecoder(bytes.NewReader(tx))
		err := decoder.Decode(&transaction)
		if err != nil {
			log.Panic(err)
		}

		txs = append(txs, transaction)
	}
	return txs
}
