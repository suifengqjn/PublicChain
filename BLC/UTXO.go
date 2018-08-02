package BLC

//Unspent Transaction output
type UTXO struct {
	// 交易ID
	TxID []byte
	// output 中的下标
	Index int
	//输出
	Output *TxOutput
}
