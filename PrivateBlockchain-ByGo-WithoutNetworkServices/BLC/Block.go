package BLC

// block本身很纯粹, 不与数据库(.db)打交道
import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	// 区块高度
	Height int64
	// 上一个区块的hash
	PrevBlockHash []byte
	// 交易数据
	Txs []*Transaction //多笔交易
	// 时间戳
	Timestamp int64
	// Hash
	Hash []byte
	// Nonce(Nonce是Number once的缩写，在密码学中Nonce是一个只被使用一次的任意或非重复的随机数值)
	Nonce int64
}

// 将Txs转化成字节数组返回
func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range block.Txs {
		txHashes = append(txHashes, tx.TxHash)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{})) //bytes.Join(,)的第二个变量是分隔符

	return txHash[:]
}

// 将区块序列化成字节数组
func (block *Block) Serialize() []byte {
	var result bytes.Buffer            // 定义缓冲区位置变量
	encoder := gob.NewEncoder(&result) // NewEncoder → 打包到result(byte.Buffer变量)地址对应位置
	err := encoder.Encode(block)       // 这步也会同时记录自定义数据类型，以及自定义数据类型的内部各变量位置
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

// 反序列化(将字节数组反序列化为区块对象)
func DeserializeBlock(blockBytes []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

// 1.创建新的区块
func NewBlock(txs []*Transaction, height int64, prevBlockhash []byte) *Block {

	// 创建新区快
	block := &Block{height, prevBlockhash, txs, time.Now().Unix(), nil, 0}
	// ↑[]byte() → 代表(将如string格式数据)强制转换为byte数组的格式
	// ↑↑time.Now()返回的是time库里的类型, 接.Unix()返回的是int64类型
	// ↑↑↑"nil"同"null"

	// 调用工作量证明的方法并且返回有效的Hash和Nonce
	pow := NewProofOfWork(block)
	hash, nonce := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	// fmt.Println()
	// 此处↑的"fmt.Println()"和"ProofOfWork.go"里面的"func...Run..."里面的"fmt.Printf("\r%x", hash)"混用, 可以显示随机hash直到hash符合目标, 然后显示下一个

	return block
}

// 2.单独写一个方法, 生成创世区块
func CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(txs, 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
