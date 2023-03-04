package BLC

import "bytes"

type TXOutput struct {
	Value         int64
	Ripemd160Hash []byte //公钥 (谁接收这笔钱)
}

// ↓↓传一个地址过来, 判断它的txInput的名字(ScripPubKey)是否和要查询的名字对应
// ↓解锁
func (txOutput *TXOutput) UnLockScripPubKeyWithAddress(address string) bool {
	// ↓ 反编码地址
	publicKeyHash := Base58Decode([]byte(address))
	hash160 := publicKeyHash[1 : len(publicKeyHash)-4]

	return bytes.Compare(txOutput.Ripemd160Hash, hash160) == 0 // ← 如果等于, 则返回真
}

func NewTxOutput(value int64, address string) *TXOutput {
	txOutput := &TXOutput{value, nil}
	txOutput.Lock(address)

	return txOutput
}

// 上锁
func (txOutput *TXOutput) Lock(address string) {
	publicKeyHash := Base58Decode([]byte(address))
	txOutput.Ripemd160Hash = publicKeyHash[1 : len(publicKeyHash)-4]
}
