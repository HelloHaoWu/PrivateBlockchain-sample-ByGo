package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

type TXOutputs struct {
	UTXOs []*UTXO
}

// 将区块序列化成字节数组
func (txOutputs *TXOutputs) Serialize() []byte {
	var result bytes.Buffer            // 定义缓冲区位置变量
	encoder := gob.NewEncoder(&result) // NewEncoder → 打包到result(byte.Buffer变量)地址对应位置
	err := encoder.Encode(txOutputs)   // 这步也会同时记录自定义数据类型，以及自定义数据类型的内部各变量位置
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

// 反序列化(将字节数组反序列化为区块对象)
func DeserializeTXOutputs(txOutputBytes []byte) *TXOutputs {
	var txOutputs TXOutputs
	decoder := gob.NewDecoder(bytes.NewReader(txOutputBytes))
	err := decoder.Decode(&txOutputs)
	if err != nil {
		log.Panic(err)
	}
	return &txOutputs
}
