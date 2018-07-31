package main

import "ketang/publicChain/BC/BLC"

func main()  {

	blockChain := BLC.CreateBlockChainWithGenesisBlock([]byte("genesis block"))

	lastBlock := blockChain.Blocks[len(blockChain.Blocks) - 1]
	blockChain.AddBlockToBlockChain([]byte("second block"),lastBlock.Hash,lastBlock.Height+1)

	lastBlock = blockChain.Blocks[len(blockChain.Blocks) - 1]
	blockChain.AddBlockToBlockChain([]byte("second block"),lastBlock.Hash,lastBlock.Height+1)

	lastBlock = blockChain.Blocks[len(blockChain.Blocks) - 1]
	blockChain.AddBlockToBlockChain([]byte("second block"),lastBlock.Hash,lastBlock.Height+1)
}

