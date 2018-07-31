package BLC

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
)


// hash 256 bit
type Block struct{
	//字段属性
	//1.高度：区块在区块链中的编号，第一个区块页叫创世区块，为0
	Height uint64
	//2.上一个区块的Hash值
	PrevBlockHash []byte
	//3.数据：data，交易数据
	Data []byte
	//4.时间戳
	TimeStamp uint64
	//5.自己的hash
	Hash []byte

	//6.Nonce
	Nonce uint64
}



//创建一个区块链，包含一个创世区块
func CreateGenesisBlock(data []byte) *Block  {
	return NewBlock(0,data,[]byte{})
}

func NewBlock(height uint64,data []byte, preHash []byte) *Block  {
	//创建区块
	block := &Block{Height:height,
	PrevBlockHash:preHash,
	Data:data,
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

// 反序列化
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
