package BLC

import (
	"encoding/hex"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
)

//UTXO 模型
/*
UTXO分为两部分 input 和 output
input：币从哪里来（系统产生，或者发起转账者）
output：币到哪里去
*/
type Transaction struct {
	//交易id -->每一笔交易的hash
	TxID []byte
	//输入
	Vins []*TxInput

	//输出
	Vouts []*TxOutput
}


/*
交易：
1.CoinBase交易：创世区块中
2.转账产生的普通交易：
 */

func NewCoinBaseTransaction(address string) *Transaction {
	txInput := &TxInput{[]byte{}, -1, "Genesis Data"}
	txOutput := &TxOutput{CoinbaseReward, address}
	txCoinBaseTransaction := &Transaction{[]byte{}, []*TxInput{txInput}, []*TxOutput{txOutput}}
	//设置交易ID
	txCoinBaseTransaction.SetID()
	return txCoinBaseTransaction
}

//交易ID--->根据tx，生成一个hash
func (tx *Transaction) SetID() {
	//1.tx--->[]byte
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	//2.[]byte-->hash
	hash := sha256.Sum256(buf.Bytes())
	//3.为tx设置ID
	tx.TxID = hash[:]
}

//根据转账的信息，创建一个普通的交易
func NewSimpleTransaction(from, to string, amount int64, bc *BlockChain, txs []*Transaction) *Transaction {
	//1.定义Input和Output的数组
	var txInputs []*TxInput
	var txOuputs [] *TxOutput

	//2.//获取本次转账要使用output
	balance, spendUtxo := bc.FindSpentAbleUTXos(from, amount, txs)

	for txID, indexArray := range spendUtxo {
		txIdBytes,_ := hex.DecodeString(txID)
		for _, index :=  range indexArray  {
			txinput := &TxInput{txIdBytes, index, from}
			txInputs = append(txInputs, txinput)
		}
	}

	//3.创建Output

	//转账
	txOutput := &TxOutput{amount, to}
	txOuputs = append(txOuputs, txOutput)

	//找零
	txOutput2 := &TxOutput{balance - amount, from}
	txOuputs = append(txOuputs, txOutput2)
	//创建交易
	tx := &Transaction{[]byte{}, txInputs, txOuputs}

	//设置交易的ID
	tx.SetID()
	return tx

}

//判断tx是否时CoinBase交易
func (tx *Transaction) IsCoinBaseTransaction() bool {

	return len(tx.Vins[0].TxID) == 0 && tx.Vins[0].Vout == -1
}
