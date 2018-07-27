package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockchainIterator struct {
	Yxh_currentHash []byte
	Yxh_db          *bolt.DB
}

func (i *BlockchainIterator) Yxh_Next() *Block {
	var block *Block

	err := i.Yxh_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.Yxh_currentHash)
		block = Yxh_DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.Yxh_currentHash = block.Yxh_PrevBlockHash

	return block
}