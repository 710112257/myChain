package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Block struct {
	Timestamp     int64 //时间戳
	Data          []byte	//数据
	PrevBlockHash []byte	//前一块的哈希
	Hash          []byte	//本块哈希
	Nonce         int	//难度值
}

//把该区块序列化成字节数组
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err!=nil{
		fmt.Errorf(err.Error())
	}
	return result.Bytes()
}

