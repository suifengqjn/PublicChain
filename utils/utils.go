package utils

import (
	"bytes"
	"encoding/binary"
	"log"
	"encoding/json"
)

func IntToHex(num uint64)[]byte{
	buff :=new(bytes.Buffer)
	//将二进制数据写入w
	//func Write(w io.Writer, order ByteOrder, data interface{}) error
	err:=binary.Write(buff,binary.BigEndian,num)
	if err!=nil{
		log.Panic(err)
	}
	//转为[]byte并返回
	return buff.Bytes()
}

/*
json解析的的函数

*/
func JSONToArray(jsonString string) []string {
	var arr [] string
	err := json.Unmarshal([]byte(jsonString), &arr)
	if err != nil {
		log.Panic(err)
	}
	return arr
}
