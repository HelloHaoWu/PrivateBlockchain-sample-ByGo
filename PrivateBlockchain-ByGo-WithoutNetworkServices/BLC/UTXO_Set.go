package BLC

import (
	"encoding/hex"
	"github.com/boltdb/bolt"
	"log"
)

// 1. 有一个方法, 功能:
// 遍历整个数据库, 读取所有的未花费的UTXO, 将所有的UTXO存储到数据库(这样查的时候不用查整个链查数据库就行)
// 有一个重置(reset)按钮, 当执行重置按钮时候, 就会遍历整个数据库, 把所有未花费交易输出导出来
// 有个方法, 能返回一个数组 → []*TXOutputs → 判断TXOutputs结构体内部的TxHash, 从而返回[]*TXOutputs

const utxoTableName = "utxoTableName"

type UTXOSet struct {
	Blockchain *Blockchain
}

// 重置数据库表
func (utxoSet *UTXOSet) ResetUTXOSet() {
	err := utxoSet.Blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			err := tx.DeleteBucket([]byte(utxoTableName)) // 要重置, 先删除老表
			if err != nil {
				log.Panic(err) // 删除失败
			}
		}
		// ↓ 创建新的表
		b, _ = tx.CreateBucket([]byte(utxoTableName))
		if b != nil {
			txOutputsMap := utxoSet.Blockchain.FindUTXOMap() // 返回一个[string]*TXOutput

			for keyHash, outs := range txOutputsMap {
				txHash, _ := hex.DecodeString(keyHash)
				b.Put(txHash, outs.Serialize())
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

func (utxoSet *UTXOSet) findUTXOForAddress(address string) []*UTXO {

	var utxos []*UTXO

	utxoSet.Blockchain.DB.View(func(tx *bolt.Tx) error {
		// 假设表存在
		b := tx.Bucket([]byte(utxoTableName))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			txOutputs := DeserializeTXOutputs(v)
			for _, utxo := range txOutputs.UTXOs {
				if utxo.Output.UnLockScripPubKeyWithAddress(address) {
					utxos = append(utxos, utxo)
				}
			}
		}
		return nil
	})
	return utxos
}

func (utxoSet *UTXOSet) GetBalance(address string) int64 {
	UTXOS := utxoSet.findUTXOForAddress(address)

	var amount int64

	for _, utxo := range UTXOS {
		amount += utxo.Output.Value
	}

	return amount
}
