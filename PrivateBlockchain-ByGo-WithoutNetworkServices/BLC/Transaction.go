package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
	"math/big"
	"time"
)

// UTXO(未花费的货币模型)
type Transaction struct {
	//1. 交易Hash
	TxHash []byte
	//2. 输入
	Vins []*TXInput
	//3. 输出
	Vouts []*TXOutput
}

// ↓判断当前这笔交易是否属于创世区块当中的transaction(或者说是否是Coinbase交易)
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return len(tx.Vins[0].Txid) == 0 && tx.Vins[0].Vout == -1 // ← 其实本质上就是在判断是不是创世区块
}

// Transaction创建分2种情况
// 1. 创世区块创建时的transaction → 创建区块时候产生的第一笔交易, 叫做coinbase transaction
func NewCoinbaseTransaction(address string) *Transaction {
	// ↓ 代表消费; 三个填入变量: 交易Hash(由于是创世区块所以无), 对应的交易转账部分的索引(由于无所以为0), 用户签名(因为是创世区块, 所以随便填)
	txInput := &TXInput{[]byte{}, -1, nil, []byte{}}
	// ↓
	txOutput := NewTxOutput(10, address)
	//txOutput := &TXOutput{10, address} // 已弃用
	// → go语言里的空是"nil"
	// ↓ 注意这里传的是TXInput和TXOutput的数组, 然后将创建信息txInput和txOutput作为数组的第一位
	txCoinbase := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}
	// 交易Hash不需要挖矿, 只是一个id而已 → 所以可以把整个transaction序列化, 再转sha256
	txCoinbase.TXHashSerialize() // ← 设置TxHash
	return txCoinbase            // ← 返回的形式: &Transaction{}
}

// ↓ 挖矿奖励代码(只奖励一个币)
func TransactionReward(address string) *Transaction {
	// ↓ 代表消费; 三个填入变量: 交易Hash(由于是创世区块所以无), 对应的交易转账部分的索引(由于无所以为0), 用户签名(因为是创世区块, 所以随便填)
	txInput := &TXInput{[]byte{}, -1, nil, []byte{}}
	// ↓
	txOutput := NewTxOutput(1, address)
	//txOutput := &TXOutput{10, address} // 已弃用
	// → go语言里的空是"nil"
	// ↓ 注意这里传的是TXInput和TXOutput的数组, 然后将创建信息txInput和txOutput作为数组的第一位
	txCoinbase := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}
	// 交易Hash不需要挖矿, 只是一个id而已 → 所以可以把整个transaction序列化, 再转sha256
	txCoinbase.TXHashSerialize() // ← 设置TxHash
	return txCoinbase            // ← 返回的形式: &Transaction{}
}

// ↓ 序列化生成transaction中的TxHash; 调用这个方法就会产生TxHash
func (tx *Transaction) TXHashSerialize() {
	var result bytes.Buffer            // 定义缓冲区位置变量
	encoder := gob.NewEncoder(&result) // NewEncoder → 打包到result(byte.Buffer变量)地址对应位置
	err := encoder.Encode(tx)          // 这步也会同时记录自定义数据类型，以及自定义数据类型的内部各变量位置
	if err != nil {
		log.Panic(err)
	}
	resultBytes := bytes.Join([][]byte{IntToHex(time.Now().Unix()), result.Bytes()}, []byte{})
	hash := sha256.Sum256(resultBytes)
	tx.TxHash = hash[:] // ←这样会直接由"*Transaction"的指针形式, 修改到传入交易的TxHash
}

// 2. 转账时产生的transaction(类比创建Coinbase交易的NewCoinbaseTransaction
func NewSimpleTransaction(from string, to string, amount int, blockchain *Blockchain, txs []*Transaction) *Transaction {
	// 1. 有一个函数, 能返回from这个人所有的未花费的交易输出所对应的Transaction
	//unUTXOs := blockchain.UnUTXOs(from) // 把它作为参数, 传入下面的获取money和dic的函数中

	// ↓ 获取原生的公钥, 创建TXInput对象要用
	wallets, _ := NewWallets()
	wallet := wallets.Wallets[from]

	// 2. 有一个函数, 能返回未花费的交易输出(即钱包存量, 这里基于目前数据库, 应该返回2); 同时返回一个dic(字典), 以便知道哪个可以支付这个输出
	money, spendableUTXODic := blockchain.FindSpendableUTXOs(from, amount, txs) // {hash1: [0, 2]}

	var txInputs []*TXInput
	var txOutputs []*TXOutput

	for txHash, indexArray := range spendableUTXODic {
		txHashBytes, _ := hex.DecodeString(txHash)
		for _, index := range indexArray {
			txInput := &TXInput{txHashBytes, index, nil, wallet.PublicKey}

			// ↓ 代表消费; 三个填入变量: 交易Hash(由于是创世区块所以无), 对应的交易转账部分的索引(由于无所以为0), 用户签名(因为是创世区块, 所以随便填); 这里直接导入from假设只记录单笔转账
			//txInput := &TXInput{txHashBytes, index, from} // 已弃用

			// ↑ 用[]byte(string), 会改变内部string的内容
			txInputs = append(txInputs, txInput) // ← append的作用是将txInput传入txInputs数组内
		}
	}

	// ↓ 转账部分
	txOutput := NewTxOutput(int64(amount), to)
	//txOutput := &TXOutput{int64(amount), to} // 已弃用
	txOutputs = append(txOutputs, txOutput) // ← append的作用是将txOutput传入txOutputs数组内

	// ↓ 找零部分
	txOutput = NewTxOutput(int64(money)-int64(amount), from)
	//txOutput = &TXOutput{int64(money) - int64(amount), from} // 已弃用
	txOutputs = append(txOutputs, txOutput)

	// → go语言里的空是"nil"
	// ↓ 注意这里传的是TXInput和TXOutput的数组, 然后将创建信息txInput和txOutput作为数组的第一位
	tx := &Transaction{[]byte{}, txInputs, txOutputs}
	// 交易Hash不需要挖矿, 只是一个id而已 → 所以可以把整个transaction序列化, 再转sha256
	tx.TXHashSerialize() // ← 设置TxHash

	// ↓ 进行数字签名
	blockchain.SignTransaction(tx, wallet.PrivateKey, txs) // ← 用私钥进行数字签名

	return tx // ← 返回的形式: &Transaction{}
}

// ↓ 交易本身序列化的方法(现在只有交易Hash序列化的方法)
func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

// ↓ 通过方法产生ID
func (tx *Transaction) Hash() []byte {
	txCopy := tx
	txCopy.TxHash = []byte{}
	hash := sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

// ↓ 签名的方法
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbaseTransaction() {
		return
	}

	for _, vin := range tx.Vins {
		if prevTXs[hex.EncodeToString(vin.Txid)].TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct.") // ← 当前的Input没有找到对应的Output
		}
	}

	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vins {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vins[inID].Signature = nil
		txCopy.Vins[inID].PublicKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inID].PublicKey = nil

		// 签名代码
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.TxHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vins[inID].Signature = signature
	}
}

// 拷贝一份新的Transaction用于签名
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []*TXInput
	var outputs []*TXOutput

	for _, vin := range tx.Vins {
		inputs = append(inputs, &TXInput{vin.Txid, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vouts {
		outputs = append(outputs, &TXOutput{vout.Value, vout.Ripemd160Hash})
	}

	txCopy := Transaction{tx.TxHash, inputs, outputs}

	return txCopy
}

// 数字签名验证
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbaseTransaction() {
		return true
	}

	for _, vin := range tx.Vins {
		if prevTXs[hex.EncodeToString(vin.Txid)].TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct.")
		}
	}

	txCopy := tx.TrimmedCopy()

	curve := elliptic.P256()

	for inID, vin := range tx.Vins {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vins[inID].Signature = nil
		txCopy.Vins[inID].PublicKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inID].PublicKey = nil

		// 私钥ID
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PublicKey)
		x.SetBytes(vin.PublicKey[:(keyLen / 2)])
		y.SetBytes(vin.PublicKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.TxHash, &r, &s) == false { // 将&rawPubKey与txCopy, &r, &s生成的东西进行比对
			return false
		}
	}

	return true
}
