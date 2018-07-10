package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"encoding/hex"
	"fmt"
	"bytes"
)

const yxh_utxoTableName  = "utxoTableName"

type UTXOSet struct {
	Yxh_Blockchain *Blockchain
}

// 重置数据库表
func (utxoSet *UTXOSet) Yxh_ResetUTXOSet()  {

	err := utxoSet.Yxh_Blockchain.Yxh_DB.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(yxh_utxoTableName))

		if b != nil {

			err := tx.DeleteBucket([]byte(yxh_utxoTableName))

			if err!= nil {
				log.Panic(err)
			}

		}

		b ,_ = tx.CreateBucket([]byte(yxh_utxoTableName))
		if b != nil {

			//[string]*TXOutputs
			txOutputsMap := utxoSet.Yxh_Blockchain.Yxh_FindUTXOMap()

			for keyHash,outs := range txOutputsMap {

				txHash,_ := hex.DecodeString(keyHash)

				b.Put(txHash,outs.Yxh_Serialize())

			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}

// 查谒地址对应的花费记录
func (utxoSet *UTXOSet) yxh_findUTXOForAddress(address string) []*UTXO{

	var utxos []*UTXO

	utxoSet.Yxh_Blockchain.Yxh_DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(yxh_utxoTableName))

		// 游标
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			txOutputs := Yxh_DeserializeTXOutputs(v)

			for _,utxo := range txOutputs.Yxh_UTXOS  {

				if utxo.Yxh_Output.Yxh_UnLockScriptPubKeyWithAddress(address) {
					utxos = append(utxos,utxo)
				}
			}
		}

		return nil
	})

	return utxos
}

// 查询余额
func (utxoSet *UTXOSet) Yxh_GetBalance(address string) int64 {

	UTXOS := utxoSet.yxh_findUTXOForAddress(address)

	var amount int64

	for _,utxo := range UTXOS  {
		amount += utxo.Yxh_Output.Yxh_Value
	}

	return amount
}

// 返回要凑多少钱，对应TXOutput的TX的Hash和index
func (utxoSet *UTXOSet) Yxh_FindUnPackageSpendableUTXOS(from string, txs []*Transaction) []*UTXO {

	var unUTXOs []*UTXO

	spentTXOutputs := make(map[string][]int)

	for _,tx := range txs {

		if tx.Yxh_IsCoinbaseTransaction() == false {

			for _, in := range tx.Yxh_Vins {
				//是否能够解锁
				publicKeyHash := Yxh_Base58Decode([]byte(from))

				ripemd160Hash := publicKeyHash[1:len(publicKeyHash) - 4]

				if in.Yxh_UnLockRipemd160Hash(ripemd160Hash) {

					key := hex.EncodeToString(in.Yxh_TxHash)

					spentTXOutputs[key] = append(spentTXOutputs[key], in.Yxh_Vout)
				}

			}
		}
	}

	for _,tx := range txs {

	Work1:
		for index,out := range tx.Yxh_Vouts {

			if out.Yxh_UnLockScriptPubKeyWithAddress(from) {

				fmt.Println(from)

				fmt.Println(spentTXOutputs)

				if len(spentTXOutputs) == 0 {
					utxo := &UTXO{tx.Yxh_TxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash,indexArray := range spentTXOutputs {

						txHashStr := hex.EncodeToString(tx.Yxh_TxHash)

						if hash == txHashStr {

							var isUnSpentUTXO bool

							for _,outIndex := range indexArray {

								if index == outIndex {
									isUnSpentUTXO = true
									continue Work1
								}

								if isUnSpentUTXO == false {
									utxo := &UTXO{tx.Yxh_TxHash, index, out}
									unUTXOs = append(unUTXOs, utxo)
								}
							}
						} else {
							utxo := &UTXO{tx.Yxh_TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				}

			}

		}

	}

	return unUTXOs

}

// 查询花费支出
func (utxoSet *UTXOSet) Yxh_FindSpendableUTXOS(from string,amount int64,txs []*Transaction) (int64,map[string][]int)  {

	unPackageUTXOS := utxoSet.Yxh_FindUnPackageSpendableUTXOS(from,txs)

	spentableUTXO := make(map[string][]int)

	var money int64 = 0

	for _,UTXO := range unPackageUTXOS {

		money += UTXO.Yxh_Output.Yxh_Value;
		txHash := hex.EncodeToString(UTXO.Yxh_TxHash)
		spentableUTXO[txHash] = append(spentableUTXO[txHash],UTXO.Yxh_Index)
		if money >= amount{
			return  money,spentableUTXO
		}
	}

	// 钱还不够
	utxoSet.Yxh_Blockchain.Yxh_DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(yxh_utxoTableName))

		if b != nil {

			c := b.Cursor()
			UTXOBREAK:
			for k, v := c.First(); k != nil; k, v = c.Next() {

				txOutputs := Yxh_DeserializeTXOutputs(v)

				for _,utxo := range txOutputs.Yxh_UTXOS {

					money += utxo.Yxh_Output.Yxh_Value
					txHash := hex.EncodeToString(utxo.Yxh_TxHash)
					spentableUTXO[txHash] = append(spentableUTXO[txHash],utxo.Yxh_Index)

					if money >= amount {
						 break UTXOBREAK;
					}
				}
			}

		}

		return nil
	})

	if money < amount{
		log.Panic("余额不足......")
	}

	return  money,spentableUTXO
}


// 更新
func (utxoSet *UTXOSet) Yxh_Update()  {

	// 最新的Block
	block := utxoSet.Yxh_Blockchain.Yxh_Iterator().Yxh_Next()

	ins := []*TXInput{}

	outsMap := make(map[string]*TXOutputs)

	// 找到所有我要删除的数据
	for _,tx := range block.Yxh_Txs {

		for _,in := range tx.Yxh_Vins {
			ins = append(ins,in)
		}
	}

	for _,tx := range block.Yxh_Txs  {

		utxos := []*UTXO{}

		for index,out := range tx.Yxh_Vouts  {

			isSpent := false

			for _,in := range ins  {

				if in.Yxh_Vout == index && bytes.Compare(tx.Yxh_TxHash ,in.Yxh_TxHash) == 0 && bytes.Compare(out.Yxh_Ripemd160Hash,Yxh_Ripemd160Hash(in.Yxh_PublicKey)) == 0 {

					isSpent = true
					continue
				}
			}

			if isSpent == false {
				utxo := &UTXO{tx.Yxh_TxHash,index,out}
				utxos = append(utxos,utxo)
			}

		}

		if len(utxos) > 0 {
			txHash := hex.EncodeToString(tx.Yxh_TxHash)
			outsMap[txHash] = &TXOutputs{utxos}
		}

	}

	err := utxoSet.Yxh_Blockchain.Yxh_DB.Update(func(tx *bolt.Tx) error{

		b := tx.Bucket([]byte(yxh_utxoTableName))

		if b != nil {
			// 删除
			for _,in := range ins {

				txOutputsBytes := b.Get(in.Yxh_TxHash)

				if len(txOutputsBytes) == 0 {
					continue
				}

				fmt.Println(txOutputsBytes)

				txOutputs := Yxh_DeserializeTXOutputs(txOutputsBytes)

				fmt.Println(txOutputs)

				UTXOS := []*UTXO{}

				// 判断是否需要
				isNeedDelete := false

				for _,utxo := range txOutputs.Yxh_UTXOS  {

					if in.Yxh_Vout == utxo.Yxh_Index && bytes.Compare(utxo.Yxh_Output.Yxh_Ripemd160Hash,Yxh_Ripemd160Hash(in.Yxh_PublicKey)) == 0 {

						isNeedDelete = true
					} else {
						UTXOS = append(UTXOS,utxo)
					}
				}

				if isNeedDelete {
					b.Delete(in.Yxh_TxHash)
					if len(UTXOS) > 0 {

						preTXOutputs := outsMap[hex.EncodeToString(in.Yxh_TxHash)]

						preTXOutputs.Yxh_UTXOS = append(preTXOutputs.Yxh_UTXOS,UTXOS...)

						outsMap[hex.EncodeToString(in.Yxh_TxHash)] = preTXOutputs

					}
				}

			}

			// 新增
			for keyHash,outPuts := range outsMap  {
				keyHashBytes,_ := hex.DecodeString(keyHash)
				b.Put(keyHashBytes,outPuts.Yxh_Serialize())
			}

		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}




