package BLC

import "math/big"

//区块链
type BlockChain struct {
	Blocks []*Block
}

//创建一个区块链，包含创世区块
func CreateBlockChainWithGenesisBlock(data []byte) *BlockChain {
	//1.创建创世区块
	genesisBlock := CreateGenesisBlock(data)
	//2.创建区块链对象并返回
	return &BlockChain{[]*Block{genesisBlock}}
}

//添加区块到区块链中
func (bc *BlockChain) AddBlockToBlockChain(data []byte, prevBlockHash [] byte, height uint64) {
	//1.根据参数的数据，创建Block
	newBlock := NewBlock(height, data,prevBlockHash)
	//2.将block加入blockchain
	if bc.isValid(newBlock) {
		bc.Blocks = append(bc.Blocks, newBlock)
	}

}



//提供一个方法：当前区块hash的有效性
func (bc *BlockChain)isValid(block *Block)bool{

	lastBlock := bc.Blocks[len(bc.Blocks) - 1]
	if block.Height - 1 != lastBlock.Height {  //验证高度合法性
		return false
	}

	if block.TimeStamp <= lastBlock.TimeStamp  {  //验证时间戳合法性
		return false
	}

	hashInt :=new(big.Int)
	hashInt.SetBytes(block.Hash)

	target := big.NewInt(1)
	target = target.Lsh(target, 256-TargetBit)
	if hashInt.Cmp(target) != -1 {      //验证hash合法性
		return false
	}

	return true

}