package BLC

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
)


// hash 256 bit
type Block struct{
	//字段属性
	//1.高度：区块在区块链中的编号，第一个区块页叫创世区块，为0
	Height uint64
	//2.上一个区块的Hash值
	PrevBlockHash []byte
	//3.数据：data，交易数据
	Txs []*Transaction
	//4.时间戳
	TimeStamp uint64
	//5.自己的hash
	Hash []byte

	//6.Nonce
	Nonce uint64
}



//创建一个区块链，包含一个创世区块
func CreateGenesisBlock(txs []*Transaction) *Block  {
	return NewBlock(0,txs,[]byte{})
}

func NewBlock(height uint64,txs []*Transaction, preHash []byte) *Block  {
	//创建区块
	block := &Block{Height:height,
		PrevBlockHash:preHash,
		Txs:txs,
		TimeStamp:uint64(time.Now().Unix()),
	}
	pow := NewProofOfWork(block)
	hash,nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block

}


// 序列化block对象 返回[]byte
func (block *Block)Serialize() []byte  {

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}

// 反序列化 转化为block对象
func DeserializeBlock(blockBytes []byte) *Block  {
	var block Block
	reader := bytes.NewReader(blockBytes)
	decoder := gob.NewDecoder(reader)
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

//提供一个方法，用于将block块中的txs转为[]byte数组

func (block *Block) HashTransactions()[]byte{
	//1.创建一个二维数组，存储每笔交易的txid
	var txshashes [][] byte
	//2.遍历
	for _,tx:=range block.Txs{
		/*
		tx1,tx2,tx3...
		[][]{tx1.ID,tx2.ID,tx3.ID...}

		合并-->[]--->sha256
		 */
		txshashes  = append(txshashes,tx.TxID)
	}
	//3.生成hash
	txhash:=sha256.Sum256(bytes.Join(txshashes,[]byte{}))
	return txhash[:]
}