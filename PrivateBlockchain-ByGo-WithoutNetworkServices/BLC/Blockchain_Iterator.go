package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

// Blockchain的迭代器
type BlockchainIterator struct {
	CurrentHash []byte
	DB          *bolt.DB
}

// 进入下一个区块(一次看一个)
func (blockchainIterator *BlockchainIterator) Next() *Block {
	var block *Block
	err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockTableName))
		if b != nil {
			currentBlockBytes := b.Get(blockchainIterator.CurrentHash)
			// 获取当前迭代器中currentHash所对应的区块
			block = DeserializeBlock(currentBlockBytes)

			// 更新最新的区块
			blockchainIterator.CurrentHash = block.PrevBlockHash
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return block //返回当前Hash的区块, 但迭代器的hash已经是前一个区块的Hash
}
