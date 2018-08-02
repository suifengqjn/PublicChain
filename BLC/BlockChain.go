package BLC

import (
	"math/big"
	"github.com/boltdb/bolt"
	"os"
	"fmt"
	"log"
	"time"
	"encoding/hex"
	"strconv"
)


var blockChain *BlockChain

//区块链
type BlockChain struct {
	DB        *bolt.DB //存放区块的数据库  (hash ==> block)
	LastBlockHash []byte   //最新区块的hash
}

//创建一个区块链，包含创世区块
func CreateBlockChainWithGenesisBlock(address string) {

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

	//创建一个txs--->CoinBase
	txCoinbase := NewCoinBaseTransaction(address)
	genesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase})
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
		fmt.Println("\t交易信息：")
		for _, tx := range block.Txs {
			fmt.Printf("\t\t交易ID：%x\n", tx.TxID) //[]byte
			fmt.Println("\t\tVins:")
			for _, in := range tx.Vins { //每一个TxInput：Txid，vout，解锁脚本
				fmt.Printf("\t\t\tTxID:%x\n", in.TxID)
				fmt.Printf("\t\t\tVout:%d\n", in.Vout)
				fmt.Printf("\t\t\tScriptSiq:%s\n", in.ScriptSiq)
			}
			fmt.Println("\t\tVouts:")
			for _, out := range tx.Vouts { //每个以txOutput:value,锁定脚本
				fmt.Printf("\t\t\tValue:%d\n", out.Value)
				fmt.Printf("\t\t\tScriptPubKey:%s\n", out.ScriptPubKey)
			}
		}
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



//新增功能：通过转账，创建区块
func (bc *BlockChain) MineNewBlock(from, to, amount []string) {
	/*
	1.新建交易
	2.新建区块：
		读取数据库，获取最后一块block
	3.存入到数据库中
	 */

	//fmt.Println(from)
	//fmt.Println(to)
	//fmt.Println(amount)
	//1.新建交易

	var txs [] *Transaction

	for i := 0;i < len(from);i++ {
		//amount[0]-->int
		if v, _ :=strconv.Atoi(amount[i]); v <= 0 {
			fmt.Println(from[i], "余额不合法")
			os.Exit(1)
		}
		amountInt, _ := strconv.ParseInt(amount[i], 10, 64)
		tx := NewSimpleTransaction(from[i], to[i], amountInt, bc, txs)
		txs = append(txs, tx)
	}


	//2.新建区块
	newBlock := new(Block)
	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucketName))
		if b != nil {
			//读取数据库
			blockBytes := b.Get(bc.LastBlockHash)
			lastBlock := DeserializeBlock(blockBytes)

			newBlock = NewBlock(lastBlock.Height+1,txs, lastBlock.Hash)

		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	//3.存入到数据库中
	err = bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucketName))
		if b != nil {
			//将新block存入到数据库中
			b.Put(newBlock.Hash, newBlock.Serialize())
			//更新l
			b.Put([]byte("l"), newBlock.Hash)
			//tip
			bc.LastBlockHash = newBlock.Hash
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

//提供一个功能：查询余额
func (bc *BlockChain) GetBalance(address string) int64 {
	unSpendUTXOs := bc.UnSpent(address,[]*Transaction{})
	var total int64
	for _, utxo := range unSpendUTXOs {
		total += utxo.Output.Value
	}
	return total

}

//设计一个方法，用于获取指定用户的所有的未花费Txoutput
/*
UTXO模型：未花费的交易输出
	Unspent Transaction TxOutput
 */
func (bc *BlockChain) UnSpent(address string, txs []*Transaction) []*UTXO {//王二狗
	/*
	1.遍历数据库，获取每个block--->Txs
	2.遍历所有交易：
		Inputs，---->将数据，记录为已经花费
		Outputs,---->每个output
	 */
	//存储未花费的TxOutput
	var unSpentUTXOs []*UTXO
	//存储已经花费的信息
	spentTxOutputMap := make(map[string][]int) // map[TxID] = []int{vout}

	//第一部分：先查询本次转账，已经产生了的Transanction
	for i := len(txs)-1;i>=0;i--{
		unSpentUTXOs = caculate(txs[i],address,spentTxOutputMap,unSpentUTXOs)
	}

	it := bc.Iterator()
	//第二部分：数据库里的Trasacntion
	for {
		//1.获取每个block
		block := it.Next()
		//2.遍历该block的txs
		for i := len(block.Txs) - 1; i >= 0; i-- {
			unSpentUTXOs = caculate(block.Txs[i],address,spentTxOutputMap,unSpentUTXOs)
		}
		//3.判断推出
		hashInt := new(big.Int)
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(hashInt) == 0 {
			break
		}

	}

	return unSpentUTXOs
}



func caculate(tx *Transaction,address string, spentTxOutputMap map[string][]int,unSpentUTXOs []*UTXO) []*UTXO{
	//遍历每个tx：txID，Vins，Vouts

	//遍历所有的TxInput
	if !tx.IsCoinBaseTransaction() { //tx不是CoinBase交易，遍历TxInput
		for _, txInput := range tx.Vins {
			//txInput-->TxInput
			if txInput.UnlockWithAddress(address) {
				//txInput的解锁脚本(用户名) 如果和钥查询的余额的用户名相同，
				key := hex.EncodeToString(txInput.TxID)
				spentTxOutputMap[key] = append(spentTxOutputMap[key], txInput.Vout)
				/*
				map[key]-->value
				map[key] -->[]int
				 */
			}
		}
	}

	//遍历所有的TxOutput
outputs:
	for index, txOutput := range tx.Vouts { //index= 0,txoutput.锁定脚本：王二狗
		if txOutput.UnlockWithAddress(address) {
			if len(spentTxOutputMap) != 0 {
				var isSpentOutput bool //false
				//遍历map
				for txID, indexArray := range spentTxOutputMap { //143d,[]int{1}
					//遍历 记录已经花费的下标的数组
					for _, i := range indexArray {
						if i == index && hex.EncodeToString(tx.TxID) == txID {
							isSpentOutput = true //标记当前的txOutput是已经花费
							continue outputs
						}
					}
				}

				if !isSpentOutput {
					//unSpentTxOutput = append(unSpentTxOutput, txOutput)
					//根据未花费的output，创建utxo对象--->数组
					utxo := &UTXO{tx.TxID, index, txOutput}
					unSpentUTXOs = append(unSpentUTXOs, utxo)
				}

			} else {
				//如果map长度未0,证明还没有花费记录，output无需判断
				//unSpentTxOutput = append(unSpentTxOutput, txOutput)
				utxo := &UTXO{tx.TxID, index, txOutput}
				unSpentUTXOs = append(unSpentUTXOs, utxo)
			}
		}
	}
	return unSpentUTXOs

}



//用于一次转账的交易中，可以使用的utxo
func (bc *BlockChain) FindSpentAbleUTXos(from string, amount int64, txs []*Transaction) (int64, map[string][]int)  {
	/*
 	1.根据from获取到的所有的utxo
 	2.遍历utxos，累加余额，判断，是否如果余额，大于等于要要转账的金额，


 	返回：map[txID] -->[]int{下标1，下标2} --->Output
 	 */

 	 var balance int64
 	 spentAbleMap := make(map[string][]int)
 	 utxos := bc.UnSpent(from, txs)
 	 for _, utxo := range utxos {
 	 	balance += utxo.Output.Value
 	 	txIDStr := hex.EncodeToString(utxo.TxID)
 	 	spentAbleMap[txIDStr] = append(spentAbleMap[txIDStr], utxo.Index)
		 if balance >= amount {

		 }
	}
	if balance < amount {
		 fmt.Println(from, "余额不足")
		 os.Exit(1)
	}
	return balance, spentAbleMap
}



