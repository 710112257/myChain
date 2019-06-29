package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}
//迭代器
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

//增加新块
func (bc *Blockchain)AddBlock(data string)  {
	var lastHash []byte
	//只读类型，获取最后一个区块的哈希
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocksBucket"))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err!=nil{
		fmt.Errorf(err.Error())
	}
	//生成新块
	newBlock := NewBlock(data, lastHash)

	//挖掘到新块后，储存新块，并且更新l密钥，密钥储存新块的哈希
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocksBucket"))
		//把哈希为KEY，把区块为vlaue存入DB中
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err!=nil{
			fmt.Errorf(err.Error())
		}
		//l是个最新块哈希的索引
		err = b.Put([]byte("l"), newBlock.Hash)
		//区块链实例的最新块哈希的索引
		bc.tip = newBlock.Hash

		return nil
	})
}
//区块链迭代方法
func (bc *Blockchain) Iterator() *BlockchainIterator {
	//构建迭代器实例，最新的哈希索引和DB
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}
//迭代器只返回下一个区块

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocksBucket"))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}