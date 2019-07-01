package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

type Block struct {
	Timestamp     int64 //时间戳
	Transactions  []*Transaction	//交易
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
//得到单个块内所有的交易ID并生成唯一的散列值
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}