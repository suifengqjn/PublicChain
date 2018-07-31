package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

//区块链迭代器
type BlockChainIterator struct {
	DB *bolt.DB
	CurrentHash []byte
}

func (i *BlockChainIterator)Next() *Block  {
	var block *Block
	err := i.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BlockBucketName))
		if bucket != nil {
			encodedBlock := bucket.Get(i.CurrentHash)
			block = DeserializeBlock(encodedBlock)
			i.CurrentHash = block.PrevBlockHash

		}
		return nil

	})
	if err != nil {
		log.Panic(err)
	}

	return block
}

func (i *BlockChainIterator)BlockWithHash(hash []byte) *Block  {
	var block *Block
	err := i.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BlockBucketName))
		if bucket != nil {
			encodedBlock := bucket.Get(hash)
			block = DeserializeBlock(encodedBlock)

		}
		return nil

	})
	if err != nil {
		log.Panic(err)
	}
	return block
}