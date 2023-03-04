package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// 整个Hash是256位, 256位Hash里面前面至少要有{16}个0
const targetBit = 20 // 这里是假定挖矿难度为{16}

type ProofOfWork struct {
	Block  *Block   //当前要验证的区块
	target *big.Int //代表难度 大数存储(不是int类型) 范围-2^31~2^31-1
	// big.Int和int64的区别↓
	// int64等类型有位数限制, 而big.Int也有位数限制(但能解决int64的数据溢出问题)

}

// 挖矿用转化方式↓
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	// 将Block属性拼接成字节数组
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevBlockHash,
			pow.Block.HashTransactions(),
			IntToHex(pow.Block.Timestamp),
			IntToHex(int64(targetBit)), //不加也可以
			IntToHex(int64(nonce)),
			IntToHex(pow.Block.Height),
		}, []byte{})
	return data
}

func (proofOfWork *ProofOfWork) isvalid() bool {
	var hashInt big.Int
	hashInt.SetBytes(proofOfWork.Block.Hash)

	if proofOfWork.target.Cmp(&hashInt) == 1 {
		return true
	}
	return false
}

func (proofOfWork *ProofOfWork) Run() ([]byte, int64) {
	// ↑解释: 变量"proofOfWork"(变量proofOfWork属于ProofOfWork类型)方法 返回[]byte和int64类型
	nonce := 0
	var hashInt big.Int //存储新生成的hash值
	// ↑若为hashInt取地址(即hashInt *big.Int)(下面的proofOfWork.target.Cmp(hashInt), 则其不能为<nil>
	var hash [32]byte //全局Hash(256位, 32byte/每byte→8位)
	// ↓挖矿循环Hash的过程(POW)
	for {
		// 准备数据
		dataBytes := proofOfWork.prepareData(nonce)
		// 生成hash
		hash = sha256.Sum256(dataBytes) //默认格式[32]byte
		fmt.Printf("\r%x", hash)        // ← 这行不注释, 就会现实生成Hash的过程; 但你不注释生成过程就会很慢很卡
		// 将hash存储到hashInt
		hashInt.SetBytes(hash[:]) //hash[:] → 转化成[]byte
		// 判断hashInt是否小于Block里面的target
		// Cmp → compares x and y and returns:
		//  -1 if x < y
		//   0 if x = y
		//   1 if x > y
		if proofOfWork.target.Cmp(&hashInt) == 1 {
			break
		}
		nonce = nonce + 1
	}

	return hash[:], int64(nonce)
}

// 创建新的工作量证明对象 一个区块对应一个工作量证明对象(1 to 1)
// ↓工作量(判断区块是否有效, 无效产生新的hash) ~= 区块 + 难度(int)(难度为x, block前面有x个0)
func NewProofOfWork(block *Block) *ProofOfWork {

	// big.Int对象 (核心思路: Hash小于Target)
	// targetBit为X, 则target(256位)第i位为"1", 其他位为"0"

	// 创建一个初始值为1的target
	target := big.NewInt(1)
	// 左移256 - targetBit位置(定义难度)
	target.Lsh(target, 256-targetBit) //.Lsh是左移位方法

	return &ProofOfWork{block, target}
}
