package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"math/big"
	"time"
	"os"
	"strconv"
	"encoding/hex"
	"crypto/ecdsa"
	"bytes"
)

// 数据库名字
const yxh_dbName = "blockchain.db"

// 表的名字
const yxh_blockTableName = "blocks"


//Blockchain是链接所有区块的对象，所以关键的方法都在这里实现
type Blockchain struct {
	Yxh_Tip []byte //最新的区块的Hash
	Yxh_DB  *bolt.DB
}

// 生成迭代器
func (blockchain *Blockchain) Yxh_Iterator() *BlockchainIterator {

	return &BlockchainIterator{blockchain.Yxh_Tip, blockchain.Yxh_DB}
}

// 判断数据库是否存在
func Yxh_DBExists() bool {

	if _, err := os.Stat(yxh_dbName); os.IsNotExist(err) {
		return false
	}

	return true
}

// 遍历输出所有区块的信息
func (blc *Blockchain) Yxh_Printchain() {

	blockchainIterator := blc.Yxh_Iterator()

	for {
		block := blockchainIterator.Yxh_Next()

		fmt.Printf("Height：%d\n", block.Yxh_Height)
		fmt.Printf("PrevBlockHash：%x\n", block.Yxh_PrevBlockHash)
		fmt.Printf("Timestamp：%s\n", time.Unix(block.Yxh_Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash：%x\n", block.Yxh_Hash)
		fmt.Printf("Nonce：%d\n", block.Yxh_Nonce)
		fmt.Println("Txs:")
		for _, tx := range block.Yxh_Txs {

			fmt.Printf("%x\n", tx.Yxh_TxHash)
			fmt.Println("Vins:")
			for _, in := range tx.Yxh_Vins {
				fmt.Printf("%x\n", in.Yxh_TxHash)
				fmt.Printf("%d\n", in.Yxh_Vout)
				fmt.Printf("%x\n", in.Yxh_PublicKey)
			}

			fmt.Println("Vouts:")
			for _, out := range tx.Yxh_Vouts {
				//fmt.Println(out.Value)
				fmt.Printf("%d\n",out.Yxh_Value)
				//fmt.Println(out.Ripemd160Hash)
				fmt.Printf("%x\n",out.Yxh_Ripemd160Hash)
			}
		}

		fmt.Println("------------------------------")

		var hashInt big.Int
		hashInt.SetBytes(block.Yxh_PrevBlockHash)

		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y

		// 凉意，for循环是死循环，必须有退出条件
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break;
		}
	}

}

//// 增加区块到区块链里面
func (blc *Blockchain) Yxh_AddBlockToBlockchain(txs []*Transaction) {

	err := blc.Yxh_DB.Update(func(tx *bolt.Tx) error {

		//1. 获取表
		b := tx.Bucket([]byte(yxh_blockTableName))
		//2. 创建新区块
		if b != nil {

			// ⚠️，先获取最新区块
			blockBytes := b.Get(blc.Yxh_Tip)
			// 反序列化
			block := Yxh_DeserializeBlock(blockBytes)

			//3. 将区块序列化并且存储到数据库中
			newBlock := Yxh_NewBlock(txs, block.Yxh_Height+1, block.Yxh_Hash)
			err := b.Put(newBlock.Yxh_Hash, newBlock.Yxh_Serialize())
			if err != nil {
				log.Panic(err)
			}
			//4. 更新数据库里面"l"对应的hash
			err = b.Put([]byte("l"), newBlock.Yxh_Hash)
			if err != nil {
				log.Panic(err)
			}
			//5. 更新blockchain的Tip
			blc.Yxh_Tip = newBlock.Yxh_Hash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

//1. 创建带有创世区块的区块链
func Yxh_CreateBlockchainWithGenesisBlock(address string) *Blockchain {

	// 判断数据库是否存在
	if Yxh_DBExists() {
		fmt.Println("创世区块已经存在.......")
		os.Exit(1)
	}

	fmt.Println("正在创建创世区块.......")

	// 创建或者打开数据库
	db, err := bolt.Open(yxh_dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var genesisHash []byte

	// 关闭数据库
	err = db.Update(func(tx *bolt.Tx) error {

		// 创建数据库表
		b, err := tx.CreateBucket([]byte(yxh_blockTableName))

		if err != nil {
			log.Panic(err)
		}

		if b != nil {
			// 创建创世区块
			// 创建了一个coinbase Transaction
			txCoinbase := Yxh_NewCoinbaseTransaction(address)

			genesisBlock := Yxh_CreateGenesisBlock([]*Transaction{txCoinbase})
			// 将创世区块存储到表中
			err := b.Put(genesisBlock.Yxh_Hash, genesisBlock.Yxh_Serialize())
			if err != nil {
				log.Panic(err)
			}

			// 存储最新的区块的hash
			err = b.Put([]byte("l"), genesisBlock.Yxh_Hash)
			if err != nil {
				log.Panic(err)
			}

			genesisHash = genesisBlock.Yxh_Hash
		}

		return nil
	})

	return &Blockchain{genesisHash, db}

}

// 返回Blockchain对象,从数据库里查询对象
func Yxh_BlockchainObject() *Blockchain {

	db, err := bolt.Open(yxh_dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var tip []byte

	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(yxh_blockTableName))

		if b != nil {
			// 读取最新区块的Hash
			tip = b.Get([]byte("l"))

		}

		return nil
	})

	return &Blockchain{tip, db}
}

// 如果一个地址对应的TXOutput未花费，那么这个Transaction就应该添加到数组中返回
func (blockchain *Blockchain) Yxh_UnUTXOs(address string,txs []*Transaction) []*UTXO {

	var unUTXOs []*UTXO

	spentTXOutputs := make(map[string][]int)

	for _,tx := range txs {

		if tx.Yxh_IsCoinbaseTransaction() == false {
			for _, in := range tx.Yxh_Vins {
				//是否能够解锁
				publicKeyHash := Yxh_Base58Decode([]byte(address))

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

			if out.Yxh_UnLockScriptPubKeyWithAddress(address) {

				fmt.Println(address)

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

	blockIterator := blockchain.Yxh_Iterator()

	for {

		block := blockIterator.Yxh_Next()

		fmt.Println()

		for i := len(block.Yxh_Txs) - 1; i >= 0 ; i-- {

			tx := block.Yxh_Txs[i]

			if tx.Yxh_IsCoinbaseTransaction() == false {
				for _, in := range tx.Yxh_Vins {
					//是否能够解锁
					publicKeyHash := Yxh_Base58Decode([]byte(address))

					ripemd160Hash := publicKeyHash[1:len(publicKeyHash) - 4]

					if in.Yxh_UnLockRipemd160Hash(ripemd160Hash) {

						key := hex.EncodeToString(in.Yxh_TxHash)

						spentTXOutputs[key] = append(spentTXOutputs[key], in.Yxh_Vout)
					}

				}
			}

			// Vouts

		work:
			for index, out := range tx.Yxh_Vouts {

				if out.Yxh_UnLockScriptPubKeyWithAddress(address) {

					fmt.Println(out)
					fmt.Println(spentTXOutputs)

					if spentTXOutputs != nil {

						if len(spentTXOutputs) != 0 {

							var isSpentUTXO bool

							for txHash, indexArray := range spentTXOutputs {

								for _, i := range indexArray {
									if index == i && txHash == hex.EncodeToString(tx.Yxh_TxHash) {
										isSpentUTXO = true
										continue work
									}
								}
							}

							if isSpentUTXO == false {

								utxo := &UTXO{tx.Yxh_TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)

							}
						} else {
							utxo := &UTXO{tx.Yxh_TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}

					}
				}

			}

		}

		fmt.Println(spentTXOutputs)

		var hashInt big.Int
		hashInt.SetBytes(block.Yxh_PrevBlockHash)

		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		//退出for循环
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}

	}

	return unUTXOs
}

// 转账时查找可用的UTXO
func (blockchain *Blockchain) Yxh_FindSpendableUTXOS(from string, amount int,txs []*Transaction) (int64, map[string][]int) {

	//1. 现获取所有的UTXO

	utxos := blockchain.Yxh_UnUTXOs(from, txs)

	spendableUTXO := make(map[string][]int)

	//2. 遍历utxos

	var value int64

	for _, utxo := range utxos {

		value = value + utxo.Yxh_Output.Yxh_Value

		hash := hex.EncodeToString(utxo.Yxh_TxHash)

		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Yxh_Index)

		if value >= int64(amount) {
			break
		}
	}

	if value < int64(amount) {

		fmt.Printf("%s's fund is 不足\n", from)
		os.Exit(1)
	}

	return value, spendableUTXO
}

// 挖掘新的区块
func (blockchain *Blockchain) Yxh_MineNewBlock(from []string, to []string, amount []string) {

	//	$ ./bc send -from '["juncheng"]' -to '["zhangqiang"]' -amount '["2"]'
	//	[juncheng]
	//	[zhangqiang]
	//	[2]

	//1.建立一笔交易

	utxoSet := &UTXOSet{blockchain}

	var txs []*Transaction

	for index,address := range from {
		value, _ := strconv.Atoi(amount[index])
		tx := Yxh_NewSimpleTransaction(address, to[index], int64(value), utxoSet,txs)
		txs = append(txs, tx)
		//fmt.Println(tx)
	}

	//新区块挖矿奖励
	tx := Yxh_NewCoinbaseTransaction(from[0])
	txs = append(txs,tx)

	//1. 通过相关算法建立Transaction数组
	var block *Block

	blockchain.Yxh_DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(yxh_blockTableName))
		if b != nil {

			hash := b.Get([]byte("l"))

			blockBytes := b.Get(hash)

			block = Yxh_DeserializeBlock(blockBytes)

		}

		return nil
	})

	// 在建立新区块之前对txs进行签名验证

	_txs := []*Transaction{}

	for _,tx := range txs  {

		if blockchain.Yxh_VerifyTransaction(tx,_txs) != true {
			log.Panic("ERROR: Invalid transaction")
		}

		_txs = append(_txs,tx)
	}

	//2. 建立新的区块
	block = Yxh_NewBlock(txs, block.Yxh_Height+1, block.Yxh_Hash)

	//将新区块存储到数据库
	blockchain.Yxh_DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(yxh_blockTableName))
		if b != nil {

			b.Put(block.Yxh_Hash, block.Yxh_Serialize())

			b.Put([]byte("l"), block.Yxh_Hash)

			blockchain.Yxh_Tip = block.Yxh_Hash

		}
		return nil
	})

}

// 查询余额(可遍历区块或遍历未花费交易记录，用于SPV轻钱包查询和转帐)
func (blockchain *Blockchain) Yxh_GetBalance(address string) int64 {

	utxos := blockchain.Yxh_UnUTXOs(address,[]*Transaction{})

	var amount int64

	for _, utxo := range utxos {

		amount = amount + utxo.Yxh_Output.Yxh_Value
	}

	return amount
}

//签名
func (bclockchain *Blockchain) Yxh_SignTransaction(tx *Transaction,privKey ecdsa.PrivateKey,txs []*Transaction)  {

	//创世和挖矿产生的交易不需要签名
	if tx.Yxh_IsCoinbaseTransaction() {
		return
	}

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Yxh_Vins {
		prevTX, err := bclockchain.Yxh_FindTransaction(vin.Yxh_TxHash,txs)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.Yxh_TxHash)] = prevTX
	}

	tx.Yxh_Sign(privKey, prevTXs)
}

// 正式签名方法
func (bc *Blockchain) Yxh_FindTransaction(ID []byte,txs []*Transaction) (Transaction, error) {

	for _,tx := range txs  {
		if bytes.Compare(tx.Yxh_TxHash, ID) == 0 {
			return *tx, nil
		}
	}


	bci := bc.Yxh_Iterator()

	for {
		block := bci.Yxh_Next()

		for _, tx := range block.Yxh_Txs {
			if bytes.Compare(tx.Yxh_TxHash, ID) == 0 {
				return *tx, nil
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.Yxh_PrevBlockHash)


		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break;
		}
	}

	return Transaction{},nil
}


// 验证数字签名(挖矿前需要执行校验)
func (bc *Blockchain) Yxh_VerifyTransaction(tx *Transaction,txs []*Transaction) bool {

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Yxh_Vins {
		prevTX, err := bc.Yxh_FindTransaction(vin.Yxh_TxHash,txs)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.Yxh_TxHash)] = prevTX
	}

	return tx.Yxh_Verify(prevTXs)
}


// [string]*TXOutputs
func (blc *Blockchain) Yxh_FindUTXOMap() map[string]*TXOutputs  {

	blcIterator := blc.Yxh_Iterator()

	// 存储已花费的UTXO的信息
	spentableUTXOsMap := make(map[string][]*TXInput)

	utxoMaps := make(map[string]*TXOutputs)

	for {
		block := blcIterator.Yxh_Next()

		for i := len(block.Yxh_Txs) - 1; i >= 0 ;i-- {

			txOutputs := &TXOutputs{[]*UTXO{}}

			tx := block.Yxh_Txs[i]

			// coinbase
			if tx.Yxh_IsCoinbaseTransaction() == false {
				for _,txInput := range tx.Yxh_Vins {

					txHash := hex.EncodeToString(txInput.Yxh_TxHash)
					spentableUTXOsMap[txHash] = append(spentableUTXOsMap[txHash],txInput)

				}
			}

			txHash := hex.EncodeToString(tx.Yxh_TxHash)

			WorkOutLoop:
			for index,out := range tx.Yxh_Vouts  {

				if tx.Yxh_IsCoinbaseTransaction() {

					fmt.Println("IsCoinbaseTransaction")
					fmt.Println(out)
					fmt.Println(txHash)
				}

				txInputs := spentableUTXOsMap[txHash]

				if len(txInputs) > 0 {

					isSpent := false

					for _,in := range  txInputs {

						outPublicKey := out.Yxh_Ripemd160Hash
						inPublicKey := in.Yxh_PublicKey

						if bytes.Compare(outPublicKey,Yxh_Ripemd160Hash(inPublicKey)) == 0{
							if index == in.Yxh_Vout {
								isSpent = true
								continue WorkOutLoop
							}
						}

					}

					if isSpent == false {
						utxo := &UTXO{tx.Yxh_TxHash,index,out}
						txOutputs.Yxh_UTXOS = append(txOutputs.Yxh_UTXOS,utxo)
					}

				} else {
					utxo := &UTXO{tx.Yxh_TxHash,index,out}
					txOutputs.Yxh_UTXOS = append(txOutputs.Yxh_UTXOS,utxo)
				}

			}

			// 设置键值对
			utxoMaps[txHash] = txOutputs
		}

		// 找到创世区块时退出
		var hashInt big.Int
		hashInt.SetBytes(block.Yxh_PrevBlockHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}
	}

	return utxoMaps
}