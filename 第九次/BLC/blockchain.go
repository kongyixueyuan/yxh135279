package BLC

import (
	"github.com/boltdb/bolt"
	"os"
	"fmt"
	"log"
	"encoding/hex"
	"strconv"
	"crypto/ecdsa"
	"bytes"
	"errors"
)

const dbFile = "blockchain_%s.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "genesis data 08/07/2018 by viky"

type Blockchain struct {
	Yxh_tip []byte
	Yxh_db  *bolt.DB
}

// 打印区块链内容
func (bc *Blockchain) Yxh_Printchain() {
	bci := bc.Yxh_Iterator()

	for {
		block := bci.Yxh_Next()
		block.String()
		if len(block.Yxh_PrevBlockHash) == 0 {
			break
		}
	}

}

// 通过交易hash,查找交易
func (bc *Blockchain) Yxh_FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Yxh_Iterator()
	for {
		block := bci.Yxh_Next()
		for _, tx := range block.Yxh_Transactions {
			if bytes.Compare(tx.Yxh_ID, ID) == 0 {
				return *tx, nil
			}
		}
		if len(block.Yxh_PrevBlockHash) == 0 {
			break
		}
	}
	fmt.Printf("查找%x的交易失败",ID)
	return Transaction{}, errors.New("未找到交易")
}

// FindUTXO finds all unspent transaction outputs and returns transactions with spent outputs removed
func (bc *Blockchain) FindUTXO() map[string]TXOutputs {
	// 未花费的交易输出
	// key:交易hash   txID
	UTXO := make(map[string]TXOutputs)
	// 已经花费的交易txID : TXOutputs.index
	spentTXOs := make(map[string][]int)
	bci := bc.Yxh_Iterator()

	for {
		block := bci.Yxh_Next()

		// 循环区块中的交易
		for _, tx := range block.Yxh_Transactions {
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
		// 如果上一区块的hash为0，代表已经到创世区块，循环结束
		if len(block.Yxh_PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}

// 获取迭代器
func (bc *Blockchain) Yxh_Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.Yxh_tip, bc.Yxh_db}
}

// 新建区块链(包含创世区块)
func Yxh_CreateBlockchain(address string,nodeID string) *Blockchain {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if Yxh_dbExists(dbFile) {
		fmt.Println("blockchain数据库已经存在.")
		os.Exit(1)
	}

	var tip []byte
	cbtx := Yxh_NewCoinbaseTX(address, genesisCoinbaseData)
	genesis := Yxh_NewGenesisBlock(cbtx)

	//genesis.String()

	// 打开数据库，如果不存在自动创建
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		// 新区块存入数据库
		err = b.Put(genesis.Yxh_Hash, genesis.Yxh_Serialize())
		if err != nil {
			log.Panic(err)
		}
		// 将创世区块的hash存入数据库
		err = b.Put([]byte("l"), genesis.Yxh_Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Yxh_Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{tip, db}
}

// 获取blockchain对象
func Yxh_NewBlockchain(nodeID string) *Blockchain {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if !Yxh_dbExists(dbFile) {
		log.Panic("区块链还未创建")
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{tip, db}
}

// 生成新的区块(挖矿)
func (bc *Blockchain) MineNewBlock(from []string, to []string, amount []string,nodeID string , mineNow bool) *Block {
	UTXOSet := UTXOSet{bc}

	wallets, err := Yxh_NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}

	var txs []*Transaction

	for index, address := range from {
		value, _ := strconv.Atoi(amount[index])
		if value<=0 {
			log.Panic("错误：转账金额需要大于0")
		}
		wallet := wallets.Yxh_GetWallet(address)
		tx := Yxh_NewUTXOTransaction(&wallet, to[index], value, &UTXOSet, txs)
		txs = append(txs, tx)
	}

	if mineNow {
		// 挖矿奖励
		tx := Yxh_NewCoinbaseTX(from[0], "")
		txs = append(txs, tx)

		//=====================================
		newBlock := bc.Yxh_MineBlock(txs)
		UTXOSet.Yxh_Update(newBlock)
		return newBlock
	}else{
		// 如果不立即挖矿，将交易写到内存中
		//var txs_all []Yxh_Transaction
		//for _,value := range txs{
		//	txs_all= append(txs_all, *value)
		//}
		// 当前节点的IP地址
		nodeAddress = fmt.Sprintf("localhost:%s",nodeID)
		yxh_sendTxs(knownNodes[0],txs)
		return nil
	}


}

// 挖矿
func (bc *Blockchain) Yxh_MineBlock(txs []*Transaction) *Block  {
	var lashHash []byte
	var lastHeight int

	// 检查交易是否有效，验证签名
	for _, tx := range txs {
		if !bc.Yxh_VerifyTransaction(tx, txs) {
			log.Panic("错误：无效的交易")
		}
	}
	// 获取最后一个区块的hash,然后获取最后一个区块的信息，进而获得height
	err := bc.Yxh_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lashHash = b.Get([]byte("l"))
		blockData := b.Get(lashHash)
		block := Yxh_DeserializeBlock(blockData)
		lastHeight = block.Yxh_Height
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	// 生成新的区块
	newBlock := Yxh_NewBlock(txs, lashHash, lastHeight+1)

	// 将新区块的内容更新到数据库中
	err = bc.Yxh_db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Yxh_Hash, newBlock.Yxh_Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("l"), newBlock.Yxh_Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.Yxh_tip = newBlock.Yxh_Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	return newBlock
}

// 签名
func (bc *Blockchain) Yxh_SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey,txs []*Transaction) {
	prevTXs := make(map[string]Transaction)

	// 找到交易输入中，之前的交易
	Vin:
	for _, vin := range tx.Yxh_Vin {
		for _, tx := range txs {
			if bytes.Compare(tx.Yxh_ID, vin.Yxh_Txid) == 0 {
				prevTX := *tx
				prevTXs[hex.EncodeToString(prevTX.Yxh_ID)] = prevTX
				continue Vin
			}
		}

		prevTX, err := bc.Yxh_FindTransaction(vin.Yxh_Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.Yxh_ID)] = prevTX

	}

	tx.Yxh_Sign(privKey, prevTXs)
}

// 验证签名
func (bc *Blockchain) Yxh_VerifyTransaction(tx *Transaction,txs []*Transaction) bool {
	if tx.Yxh_IsCoinbase() {
		return true
	}

	prevTXs := make(map[string]Transaction)
	Vin:
	for _, vin := range tx.Yxh_Vin {
		for _, tx := range txs {
			if bytes.Compare(tx.Yxh_ID, vin.Yxh_Txid) == 0 {
				prevTX := *tx
				prevTXs[hex.EncodeToString(prevTX.Yxh_ID)] = prevTX
				continue Vin
			}
		}
		prevTX, err := bc.Yxh_FindTransaction(vin.Yxh_Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.Yxh_ID)] = prevTX
	}

	return tx.Yxh_Verify(prevTXs)
}

// 判断数据库是否存在
func Yxh_dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// 获取BestHeight
func (bc *Blockchain) Yxh_GetBestHeight() int {
	var lastBlock Block

	err := bc.Yxh_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash := b.Get([]byte("l"))
		blockData := b.Get(lastHash)
		lastBlock = *Yxh_DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return lastBlock.Yxh_Height
}

// 获取所有区块的hash
func (bc *Blockchain) Yxh_GetBlockHashes() [][]byte {
	var blocks [][]byte
	bci := bc.Yxh_Iterator()

	for {
		block := bci.Yxh_Next()

		blocks = append(blocks, block.Yxh_Hash)

		if len(block.Yxh_PrevBlockHash) == 0 {
			break
		}
	}

	return blocks
}

// 根据hash获取某个区块的内容
func (bc *Blockchain) Yxh_GetBlock(blockHash []byte) (Block, error) {
	var block Block

	err := bc.Yxh_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		blockData := b.Get(blockHash)

		if blockData == nil {
			return errors.New("未找到区块")
		}

		block = *Yxh_DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		return block, err
	}

	return block, nil
}

// 将区块添加到链中
func (bc *Blockchain) Yxh_AddBlock(block *Block) {
	err := bc.Yxh_db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		blockInDb := b.Get(block.Yxh_Hash)

		if blockInDb != nil {
			return nil
		}

		blockData := block.Yxh_Serialize()
		err := b.Put(block.Yxh_Hash, blockData)
		if err != nil {
			log.Panic(err)
		}

		lastHash := b.Get([]byte("l"))
		lastBlockData := b.Get(lastHash)
		lastBlock := Yxh_DeserializeBlock(lastBlockData)

		if block.Yxh_Height > lastBlock.Yxh_Height {
			err = b.Put([]byte("l"), block.Yxh_Hash)
			if err != nil {
				log.Panic(err)
			}
			bc.Yxh_tip = block.Yxh_Hash
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}