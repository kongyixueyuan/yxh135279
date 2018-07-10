package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockchainIterator struct {
	Yxh_CurrentHash []byte
	Yxh_DB  *bolt.DB
}

// 遍历一条记录
func (blockchainIterator *BlockchainIterator) Yxh_Next() *Block {

	var block *Block

	err := blockchainIterator.Yxh_DB.View(func(tx *bolt.Tx) error{

		b := tx.Bucket([]byte(yxh_blockTableName))

		if b != nil {
			currentBloclBytes := b.Get(blockchainIterator.Yxh_CurrentHash)
			//  获取到当前迭代器里面的currentHash所对应的区块
			block = Yxh_DeserializeBlock(currentBloclBytes)

			// 更新迭代器里面CurrentHash
			blockchainIterator.Yxh_CurrentHash = block.Yxh_PrevBlockHash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return block
}