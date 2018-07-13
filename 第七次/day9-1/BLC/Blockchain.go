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
const yxh_dbName = "blockchain_%s.db"

// 表的名字
const yxh_blockTableName = "blocks"

type Blockchain struct {
	Yxh_Tip []byte //最新的区块的Hash
	Yxh_DB  *bolt.DB
}

// 迭代器
func (blockchain *Blockchain) Yxh_Iterator() *BlockchainIterator {

	return &BlockchainIterator{blockchain.Yxh_Tip, blockchain.Yxh_DB}
}

// 判断数据库是否存在
//3000
//blockchain_3000.db
func Yxh_DBExists(dbName string) bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}

	return true
}

// 遍历输出所有区块的信息
func (blc *Blockchain) Yxh_Printchain() {

	fmt.Println("PrintchainPrintchainPrintchainPrintchain")
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
func Yxh_CreateBlockchainWithGenesisBlock(address string,nodeID string) *Blockchain {

	// 格式化数据库名字
	dbName := fmt.Sprintf(yxh_dbName,nodeID)


	// 判断数据库是否存在
	if Yxh_DBExists(dbName) {
		fmt.Println("创世区块已经存在.......")
		os.Exit(1)
	}

	fmt.Println("正在创建创世区块.......")

	// 创建或者打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
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
			txCoinbase := NewCoinbaseTransaction(address)

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

// 返回Blockchain对象
func Yxh_BlockchainObject(nodeID string) *Blockchain {

	dbName := fmt.Sprintf(yxh_dbName,nodeID)

	// 判断数据库是否存在
	if Yxh_DBExists(dbName) == false {
		fmt.Println("数据库不存在....")
		os.Exit(1)
	}

	db, err := bolt.Open(dbName, 0600, nil)
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

	//{hash:[0]}

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
				fmt.Println("看看是否是俊诚...")
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

		fmt.Println(block)
		fmt.Println()

		for i := len(block.Yxh_Txs) - 1; i >= 0 ; i-- {

			tx := block.Yxh_Txs[i]
			// txHash
			// Vins
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

					//&{2 zhangqiang}
					//map[]

					if spentTXOutputs != nil {

						//map[cea12d33b2e7083221bf3401764fb661fd6c34fab50f5460e77628c42ca0e92b:[0]]

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
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}

	}

	return unUTXOs
}

// 转账时查找可用的UTXO
func (blockchain *Blockchain) Yxh_FindSpendableUTXOS(from string, amount int,txs []*Transaction) (int64, map[string][]int) {

	//1. 现获取所有的UTXO

	utxos := blockchain.Yxh_UnUTXOs(from,txs)

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
func (blockchain *Blockchain) Yxh_MineNewBlock(from []string, to []string, amount []string,nodeID string) {

	//	$ ./bc send -from '["juncheng"]' -to '["zhangqiang"]' -amount '["2"]'
	//	[juncheng]
	//	[zhangqiang]
	//	[2]

	//1.建立一笔交易


	utxoSet := &UTXOSet{blockchain}

	var txs []*Transaction

	for index,address := range from {
		value, _ := strconv.Atoi(amount[index])
		tx := NewSimpleTransaction(address, to[index], int64(value), utxoSet,txs,nodeID)
		txs = append(txs, tx)
		//fmt.Println(tx)
	}

	//奖励
	tx := NewCoinbaseTransaction(from[0])
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

// 查询余额
func (blockchain *Blockchain) Yxh_GetBalance(address string) int64 {

	utxos := blockchain.Yxh_UnUTXOs(address,[]*Transaction{})

	var amount int64

	for _, utxo := range utxos {

		amount = amount + utxo.Yxh_Output.Yxh_Value
	}

	return amount
}

func (bclockchain *Blockchain) Yxh_SignTransaction(tx *Transaction,privKey ecdsa.PrivateKey,txs []*Transaction)  {

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

	tx.Sign(privKey, prevTXs)

}


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


// 验证数字签名
func (bc *Blockchain) Yxh_VerifyTransaction(tx *Transaction,txs []*Transaction) bool {


	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Yxh_Vins {
		prevTX, err := bc.Yxh_FindTransaction(vin.Yxh_TxHash,txs)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.Yxh_TxHash)] = prevTX
	}

	return tx.Verify(prevTXs)
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

			txInputs := spentableUTXOsMap[txHash]

			if len(txInputs) > 0 {


			WorkOutLoop:
				for index,out := range tx.Yxh_Vouts  {

					for _,in := range  txInputs {

						outPublicKey := out.Yxh_Ripemd160Hash
						inPublicKey := in.Yxh_PublicKey


						if bytes.Compare(outPublicKey,Yxh_Ripemd160Hash(inPublicKey)) == 0 {
							if index == in.Yxh_Vout {

								continue WorkOutLoop
							} else {

								utxo := &UTXO{tx.Yxh_TxHash,index,out}
								txOutputs.Yxh_UTXOS = append(txOutputs.Yxh_UTXOS,utxo)
							}
						}
					}


				}

			} else {

				for index,out := range tx.Yxh_Vouts {
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



//----------

func (bc *Blockchain) Yxh_GetBestHeight() int64 {

	block := bc.Yxh_Iterator().Yxh_Next()

	return block.Yxh_Height
}

func (bc *Blockchain) GetBlockHashes() [][]byte {

	blockIterator := bc.Yxh_Iterator()

	var blockHashs [][]byte

	for {
		block := blockIterator.Yxh_Next()

		blockHashs = append(blockHashs,block.Yxh_Hash)

		var hashInt big.Int
		hashInt.SetBytes(block.Yxh_PrevBlockHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}
	}

	return blockHashs
}

func (bc *Blockchain) Yxh_GetBlock(blockHash []byte) (*Block, error) {

	var block *Block

	err := bc.Yxh_DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(yxh_blockTableName))
		if b != nil {
			blockBytes := b.Get(blockHash)
			block = Yxh_DeserializeBlock(blockBytes)
		}
		return nil
	})

	return block, err
}

func (bc *Blockchain) addBlock(block *Block) {

	err := bc.Yxh_DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(yxh_blockTableName))
		if b != nil {

			blockExist := b.Get([]byte(block.Yxh_Hash))

			if blockExist != nil {
				return nil
			}

			err := b.Put(block.Yxh_Hash, block.Yxh_Serialize())

			if err != nil {
				log.Panic(err)
			}

			blockHash := b.Get([]byte("l"))
			blockBytes := b.Get(blockHash)
			blockDB := Yxh_DeserializeBlock(blockBytes)

			//增加新区块后要更新"l"记录和blickchain中的Tip值
			if blockDB.Yxh_Height < block.Yxh_Height {
				b.Put([]byte("l"), block.Yxh_Hash)
				bc.Yxh_Tip = block.Yxh_Hash
			}

		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}


}