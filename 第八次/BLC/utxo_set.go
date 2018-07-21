package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"encoding/hex"
	"fmt"
	"strings"
)

const utxoBucket = "chainstate"

type UTXOSet struct {
	Yxh_Blockchain *Blockchain
}

// 查询可花费的交易输出
func (u UTXOSet) Yxh_FindSpendableOutputs(pubkeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := u.Yxh_Blockchain.Yxh_db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := Yxh_DeserializeOutputs(v)

			for outIdx, out := range outs.Yxh_Outputs {
				if out.Yxh_IsLockedWithKey(pubkeyHash) && accumulated < amount {
					accumulated += out.Yxh_Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return accumulated, unspentOutputs
}

func (u UTXOSet) Yxh_Reindex() {
	db := u.Yxh_Blockchain.Yxh_db
	bucketName := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		// 删除旧的bucket
		err := tx.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			log.Panic()
		}
		_, err = tx.CreateBucket(bucketName)
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	UTXO := u.Yxh_Blockchain.FindUTXO()

	err = db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucketName)

		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(key, outs.Yxh_Serialize())
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
}

// 生成新区块的时候，更新UTXO数据库
func (u UTXOSet) Update(block *Block) {
	err := u.Yxh_Blockchain.Yxh_db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		for _, tx := range block.Yxh_Transactions {
			if !tx.Yxh_IsCoinbase() {
				for _, vin := range tx.Yxh_Vin {
					updatedOuts := TXOutputs{}
					outsBytes := b.Get(vin.Yxh_Txid)
					outs := Yxh_DeserializeOutputs(outsBytes)

					// 找出Vin对应的outputs,过滤掉花费的
					for outIndex, out := range outs.Yxh_Outputs {
						if outIndex != vin.Yxh_Vout {
							updatedOuts.Yxh_Outputs = append(updatedOuts.Yxh_Outputs, out)
						}
					}
					// 未花费的交易输出TXOutput为0
					if len(updatedOuts.Yxh_Outputs) == 0 {
						err := b.Delete(vin.Yxh_Txid)
						if err != nil {
							log.Panic(err)
						}
					} else { // 未花费的交易输出TXOutput>0
						err := b.Put(vin.Yxh_Txid, updatedOuts.Yxh_Serialize())
						if err != nil {
							log.Panic(err)
						}
					}
				}
			}

			// 将所有的交易输出TXOutput存入数据库中
			newOutputs := TXOutputs{}
			for _, out := range tx.Yxh_Vout {
				newOutputs.Yxh_Outputs = append(newOutputs.Yxh_Outputs, out)
			}
			err := b.Put(tx.Yxh_ID, newOutputs.Yxh_Serialize())
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// 打出某个公钥hash对应的所有未花费输出
func (u *UTXOSet) Yxh_FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput

	err := u.Yxh_Blockchain.Yxh_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := Yxh_DeserializeOutputs(v)

			for _, out := range outs.Yxh_Outputs {
				if out.Yxh_IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return UTXOs
}

// 查询某个地址的余额
func (u *UTXOSet) Yxh_GetBalance(address string) int {
	balance := 0
	pubKeyHash := Yxh_Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := u.Yxh_FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Yxh_Value
	}
	return balance
}

// 打印所有的UTXO
func (u *UTXOSet) String() {
	//outputs := make(map[string][]Yxh_TXOutput)

	var lines []string
	lines = append(lines, "---ALL UTXO:")
	err := u.Yxh_Blockchain.Yxh_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := Yxh_DeserializeOutputs(v)

			lines = append(lines, fmt.Sprintf("     Key: %s", txID))
			for i, out := range outs.Yxh_Outputs {
				//outputs[txID] = append(outputs[txID], out)
				lines = append(lines, fmt.Sprintf("     Output: %d", i))
				lines = append(lines, fmt.Sprintf("         value:  %d", out.Yxh_Value))
				lines = append(lines, fmt.Sprintf("         PubKeyHash:  %x", out.Yxh_PubKeyHash))
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(strings.Join(lines, "\n"))
}
