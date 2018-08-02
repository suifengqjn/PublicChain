package BLC

import (
	"math/big"
	"bytes"
	"ketang/publicChain/BC/utils"
	"crypto/sha256"
	"fmt"
)

const TargetBit = 16 //目标hash0的个数


type ProofOfWork struct {
	Block *Block      //要验证的block
	Target *big.Int   //目标hash
}

//自动调整难度的target计算
//按10分钟一个区块生成速度，2016个区块生成时间为2016*10分钟=14天。
//新目标值= 当前目标值 * 实际2016个区块出块时间 / 理论2016个区块出块时间(2周)。

//判断是否需要更新目标值( 2016的整数倍)，如果不是则继续使用最后一个区块的目标值
//计算前2016个区块出块用时
//如果用时低于半周，则按半周计算。防止难度增加4倍以上。
//如果用时高于8周，则按8周计算。防止难度降低到4倍以下。
//用时乘以当前难度
//再除以2周
//如果超过最大难度限制，则按最大难度处理

/*
var (
	bigOne = big.NewInt(1)
	// 最大难度：00000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffff，2^224，0x1d00ffff
	mainPowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 224), bigOne)
	powTargetTimespan = time.Hour * 24 * 14 // 两周
)
func CalculateNextWorkTarget(prev2016block, lastBlock Block) *big.Int {
	// 如果新区块(+1)不是2016的整数倍，则不需要更新，仍然是最后一个区块的 bits
	if (lastBlock.Head.Height+1)%2016 != 0 {
		return CompactToBig(lastBlock.Head.Bits)
	}
	// 计算 2016个区块出块时间
	actualTimespan := lastBlock.Head.Timestamp.Sub(prev2016block.Head.Timestamp)
	if actualTimespan < powTargetTimespan/4 {
		actualTimespan = powTargetTimespan / 4
	} else if actualTimespan > powTargetTimespan*4 {
		// 如果超过8周，则按8周计算
		actualTimespan = powTargetTimespan * 4
	}
	lastTarget := CompactToBig(lastBlock.Head.Bits)
	// 计算公式： target = lastTarget * actualTime / expectTime
	newTarget := new(big.Int).Mul(lastTarget, big.NewInt(int64(actualTimespan.Seconds())))
	newTarget.Div(newTarget, big.NewInt(int64(powTargetTimespan.Seconds())))
	//超过最多难度，则重置
	if newTarget.Cmp(mainPowLimit) > 0 {
		newTarget.Set(mainPowLimit)
	}
	return newTarget
}

*/

func NewProofOfWork(block *Block) *ProofOfWork  {
	//
	pow := &ProofOfWork{}
	pow.Block = block

	/*
	hash: 256bit
	16进制 4个0 => 2进制  16个0

	*/

	target := big.NewInt(1)
	target = target.Lsh(target, 256-TargetBit)
	pow.Target = target  //左移256-16

	return pow

}

//计算有效的hash
func (pow *ProofOfWork)Run() ([]byte, uint64)  {

	var nonce uint64 = 0
	var hash [32]byte
	for {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%d,%x",nonce,hash)
		hashInt := new(big.Int)
		hashInt = hashInt.SetBytes(hash[:])

		//-1 if x <  y
		//0 if x == y
		//+1 if x >  y

		if hashInt.Cmp(pow.Target) == -1  {
			break
		}

		nonce++

	}
	fmt.Println()
	return hash[:], nonce
}

//根据nonce，获取pow中要验证的block拼接成的数组的数据
func (pow *ProofOfWork) prepareData(nonce uint64) []byte {
	//1.根据nonce，生成pow中要验证的block的数组
	data := bytes.Join([][]byte{
		utils.IntToHex(pow.Block.Height),
		pow.Block.PrevBlockHash,
		utils.IntToHex(pow.Block.TimeStamp),
		pow.Block.HashTransactions(),
		utils.IntToHex(nonce),
		utils.IntToHex(TargetBit),
	}, []byte{})
	return data

}


//提供一个方法：当前区块hash的有效性
func (pow *ProofOfWork) IsValid()bool{
	hashInt :=new(big.Int)
	hashInt.SetBytes(pow.Block.Hash)
	return pow.Target.Cmp(hashInt) == 1
}
