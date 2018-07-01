package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

//遍历区块的跌代器
type BlockchainIterator struct {
	CurrentHash []byte
	DB  *bolt.DB
}

//遍历区块
func (blockchainIterator *BlockchainIterator) Next() *Block {

	var block *Block
	//查询数据
	err := blockchainIterator.DB.View(func(tx *bolt.Tx) error{
		//查询表
		b := tx.Bucket([]byte(blockTableName))

		if b != nil {
			currentBloclBytes := b.Get(blockchainIterator.CurrentHash)
			//  获取到当前迭代器里面的currentHash所对应的区块
			block = DeserializeBlock(currentBloclBytes)

			// 更新迭代器里面CurrentHash
			blockchainIterator.CurrentHash = block.PrevBlockHash
		}
		//数据库中最后返回 nil表示沅异常可以提交
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return block

}