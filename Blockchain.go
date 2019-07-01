package main

import (
	"encoding/hex"
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
func (bc *Blockchain)AddBlock(transactions []*Transaction)  {
	var lastHash []byte
	//只读类型，获取最后一个区块的哈希
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocksBucket"))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	//生成新块
	newBlock := NewBlock(transactions, lastHash)
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
//找到包含未使用的输出
func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	//生成迭代器
	bci := bc.Iterator()
	//迭代查找
	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			//解码交易的哈希ID
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// 把Vout转化成spentTXOs的格式
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				//判断该输出是否该地址可以使用
				if out.CanBeUnlockedWith(address) {
					//未花费的集合
					unspentTXs = append(unspentTXs, *tx)
				}
			}
			//判断是否奖励币
			if tx.IsCoinbase() == false {
				//如果不是奖励币
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						//花费的集合
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs
}
//查找未花费的输出
func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	//找到未花费的输出的集合
	unspentTransactions := bc.FindUnspentTransactions(address)

	//获得能够被该地址解锁的的未花费集合
	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}
//找到可花费的输出
func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	//找到包含未使用的输出
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {	//遍历未使用的输出集合
		//输出的ID
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {	//遍历所有输出
			//判断是否可用的解锁脚本，并且金额不低于0
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				//该用户能够使用的未使用的输出
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}
