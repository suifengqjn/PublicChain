package BLC

import (
	"math/big"
	"github.com/boltdb/bolt"
	"os"
	"fmt"
	"log"
	"time"
)


var blockChain *BlockChain

//区块链
type BlockChain struct {
	DB        *bolt.DB //存放区块的数据库  (hash ==> block)
	LastBlockHash []byte   //最新区块的hash
}

//创建一个区块链，包含创世区块
func CreateBlockChainWithGenesisBlock(data []byte) {

	/*
	1.判断数据库如果存在，直接结束方法
	2.数据库不存在，创建创世区块，并存入到数据库中
	 */
	if dbExists(){
		fmt.Println("数据库已经存在，无法创建创世区块。。")
		return
	}

	//数据库不存在
	fmt.Println("数据库不存在。。")
	fmt.Println("正在创建创世区块。。。。。")
	/*
	1.创建创世区块
	2.存入到数据库中
	 */
	genesisBlock := CreateGenesisBlock(data)
	db, err := bolt.Open(DBName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		//创世区块序列化后，存入到数据库中
		b, err := tx.CreateBucketIfNotExists([]byte(BlockBucketName))
		if err != nil {
			log.Panic(err)
		}

		if b != nil {
			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			b.Put([]byte("l"), genesisBlock.Hash)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}


}


//添加区块到区块链中
func (bc *BlockChain) AddBlockToBlockChain(data []byte) {
	//1.根据参数的数据，创建Block
	//newBlock := NewBlock(data, prevBlockHash, height)
	//2.将block加入blockchain
	//bc.Blocks = append(bc.Blocks, newBlock)
	/*
	1.操作bc对象，获取DB
	2.创建新的区块
	3.序列化后存入到数据库中
	 */
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		//打开bucket
		b := tx.Bucket([]byte(BlockBucketName))
		if b != nil {
			//获取bc的Tip就是最新hash，从数据库中读取最后一个block：hash，height
			blockByets := b.Get(bc.LastBlockHash)
			lastBlock := DeserializeBlock(blockByets) //数据库中的最后一个区块
			//创建新的区块
			newBlock := NewBlock(lastBlock.Height + 1, data, lastBlock.Hash)
			//序列化后存入到数据库中
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			//更新：bc的lastBlock，以及数据库中l的值
			b.Put([]byte("l"), newBlock.Hash)
			bc.LastBlockHash = newBlock.Hash

		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}


//提供一个方法：当前区块hash的有效性
func (bc *BlockChain) isValid(block *Block) bool {

	//if block.Height-1 != blockChain.LastBlock.Height { //验证高度合法性
	//	return false
	//}
	//
	//if block.TimeStamp <= blockChain.LastBlock.TimeStamp { //验证时间戳合法性
	//	return false
	//}

	hashInt := new(big.Int)
	hashInt.SetBytes(block.Hash)

	target := big.NewInt(1)
	target = target.Lsh(target, 256-TargetBit)
	if hashInt.Cmp(target) != -1 { //验证hash合法性
		return false
	}

	return true

}

//判断数据库文件是否存在
func dbExists() bool {
	if _, err := os.Stat(DBName); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

//判断bucket是否存在
func bucketExits(DB *bolt.DB) bool  {

	var exits bool
	DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BlockBucketName))
		if bucket != nil {
			exits = true
		} else {
			exits =false
		}
		return nil
	})
	return exits
}


func (bc *BlockChain) PrintChains() {
	/*
	.bc.DB.View(),
		根据hash，获取block的数据
		反序列化
		打印输出


	 */

	//获取迭代器
	it := bc.Iterator()
	for {
		//step1：根据currenthash获取对应的区块
		block := it.Next()
		fmt.Printf("第%d个区块的信息：\n", block.Height+1)
		fmt.Printf("\t高度：%d\n", block.Height)
		fmt.Printf("\t上一个区块Hash：%x\n", block.PrevBlockHash)
		fmt.Printf("\t自己的Hash：%x\n", block.Hash)
		fmt.Printf("\t数据：%s\n", block.Data)
		fmt.Printf("\t随机数：%d\n", block.Nonce)
		//fmt.Printf("\t时间：%d\n", block.TimeStamp)
		fmt.Printf("\t时间：%s\n", time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")) // 时间戳-->time-->Format("")

		//step2：判断block的prevBlcokhash为0,表示该block是创世取块，将结束循环
		hashInt := new(big.Int)
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(hashInt) == 0 {
			/*
			x.Cmp(y)
				-1 x < y
				0 x = y
				1 x > y
			 */
			break
		}

	}
}

//获取blockchainiterator的对象
func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{bc.DB, bc.LastBlockHash}
}


//提供一个函数，专门用于获取BlockChain对象
func GetBlockChainObject() *BlockChain{
	/*
		1.数据库存在，读取数据库，返回blockchain即可
		2.数据库 不存储，返回nil
	 */

	if dbExists() {
		//fmt.Println("数据库已经存在。。。")
		//打开数据库
		db, err := bolt.Open(DBName, 0600, nil)
		if err != nil {
			log.Panic(err)
		}

		var blockchain *BlockChain

		err = db.View(func(tx *bolt.Tx) error {
			//打开bucket，读取l对应的最新的hash
			b := tx.Bucket([]byte(BlockBucketName))
			if b != nil {
				//读取最新hash
				hash := b.Get([]byte("l"))
				blockchain = &BlockChain{db, hash}
			}
			return nil
		})
		if err != nil {
			log.Panic(err)
		}
		return blockchain
	}else{
		fmt.Println("数据库不存在，无法获取BlockChain对象。。。")
		return  nil
	}
}