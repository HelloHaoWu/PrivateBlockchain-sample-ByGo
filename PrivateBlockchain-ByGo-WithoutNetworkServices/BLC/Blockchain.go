package BLC

// blockchain需要与数据库(.db)打交道
import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"
)

// 数据库名称
const dbName = "blockchain.db"

// 数据库中表的名称
const BlockTableName = "blocks"

// Blockchain更多地反映地是一个状态↓
type Blockchain struct {
	Tip []byte //最新的区块的Hash → 仅保存唯一Hash
	DB  *bolt.DB
}

// Blockchain转BlockchainIterator(返回Blockchain的迭代器)
func (blockhain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blockhain.Tip, blockhain.DB}
}

// 判断数据库是否存在
func dbExists() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}

// 遍历区块链(一次性遍历)
func (blc *Blockchain) Printchain() {
	blockchainIterator := blc.Iterator()
	for {
		block := blockchainIterator.Next()
		fmt.Printf("Height: %d\n", block.Height)
		fmt.Printf("PrevBlockHash: %v\n", block.PrevBlockHash)
		fmt.Printf("Timestamp: %s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash: %v\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Println("Txs:") // ←这个不用加"\n"也换行
		for _, tx := range block.Txs {
			fmt.Printf("TxHash: %x\n", tx.TxHash) // []byte的返回值显示用%x
			fmt.Println("Vins:")                  // ←这个不用加"\n"也换行
			for _, in := range tx.Vins {
				fmt.Printf("transfer: %x\n", in.Txid)    // 交易中的某一笔转账id; []byte的返回值显示用%x; 为空则不显示
				fmt.Printf("index: %d\n", in.Vout)       // 在Vouts(转账详情页)中的索引
				fmt.Printf("pubkey: %v\n", in.PublicKey) // 用户名字
			}
			fmt.Println("Vouts:") // ←这个不用加"\n"也换行
			for _, out := range tx.Vouts {
				fmt.Println("Vout:")
				fmt.Printf("amount: %d\n", out.Value)                // 转账金额
				fmt.Printf("Ripemd160Hash: %x\n", out.Ripemd160Hash) // 转账人员名称
			}
		}
		fmt.Println("---------------------------------------")
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
}

// 增加区块到区块链|↓给某个类型写子函数, 其中第一个变量为该方法操作的变量名, 第二个变量为该方法返回结果
func (blc *Blockchain) AddBlockToBlockchain(txs []*Transaction) {
	// 往链里面添加区块
	err := blc.DB.Update(func(tx *bolt.Tx) error {
		// 1.获取表
		b := tx.Bucket([]byte(BlockTableName))
		// 2.向数据库中添加新区快
		if b != nil {
			// 获取最新区快的内容
			blockBytes := b.Get(blc.Tip)
			// 反序列化
			block := DeserializeBlock(blockBytes)

			// 2.1 创建新区快
			newBlock := NewBlock(txs, block.Height, block.Hash) // 变量数量必须完全符合条件
			// 2.2 将区块序列化并且存储到数据库中
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			// 2.3 更新数据库中的"list"
			err = b.Put([]byte("list"), newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			// 2.4 更新blockchain的Tip
			blc.Tip = newBlock.Hash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

// 1.创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock(address string) *Blockchain {
	// 判断数据库是否存在
	if dbExists() {
		fmt.Println("创世区块已经存在.")
		os.Exit(1)
	}

	fmt.Println("正在创建创世区块...")

	// 创建或打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var genesisHash []byte
	//// 这里之所以定义blockHash, 是因为genesisBlock是函数内变量, 取不到函数外
	//var blockHash []byte

	// 创建表/读取表(创建创世区块, 表一定不存在)
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(BlockTableName))
		if err != nil {
			log.Panic(err)
		}
		// 这里的思路是将区块链创建和数据库建立并行操作↓
		if b != nil {
			// 创建创世区块
			// ↓创建一个coinbase Transaction
			txCoinbase := NewCoinbaseTransaction(address)                  // ←这里相比于原来, 多了一步由地址(address)生成原始交易(transaction)的过程
			genesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase}) // ←创建创世区块的过程, 需要一系列交易; 这里把一个包含第一笔原始交易的交易数组返回
			// 将创世区块存储到表中
			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			// 存储最新的区块的hash
			err = b.Put([]byte("list"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			//// genesisBlock的Hash要通过这个该函数的"全局变量"取出来
			//blockHash = genesisBlock.Hash

			genesisHash = genesisBlock.Hash // ← 声明过就不用写:=了
		}
		return nil
	})

	return &Blockchain{genesisHash, db}

	//// 返回区块链对象
	//return &Blockchain{blockHash, db}
}

// 返回Blockchain对象
func GetBlockchainObject() *Blockchain {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockTableName))
		if b != nil {
			// 读取最新区块的Hash
			tip = b.Get([]byte("list"))
		}
		return nil
	})
	return &Blockchain{tip, db}
}

// ↓如果一个地址对应的TxOutput未花费, 那么这个Transaction就应该添加到数组中返回 →→ 这里修改成把所有未花费的Output返回
func (blockchain *Blockchain) UnUTXOs(address string, txs []*Transaction) []*UTXO {
	// 需要做的: 遍历数据库
	// ↓用于存储找到的transaction
	var unUTXOs []*UTXO
	// ↓已花费的输出(字典格式)
	spentTXOutputs := make(map[string][]int) // ← 格式: hash → 索引
	for _, tx := range txs {
		if tx.IsCoinbaseTransaction() == false { // ← 这行代码说明其是其他正常的transaction
			for _, in := range tx.Vins {
				publicKeyHash := Base58Decode([]byte(address))           // 反编码
				ripemd160Hash := publicKeyHash[1 : len(publicKeyHash)-4] // ← 取公钥Hash部分
				// ↓是否能够解锁
				if in.UnLockRipemd160Hash(ripemd160Hash) {
					key := hex.EncodeToString(in.Txid)
					spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout) // ← 这个"Vout"不是"Vouts", 而是TXInputs里面标注的其继承的交易的索引
					// ↑不用获取全部的Inputs后再核验Output, 因为每个Input只有可能调用自之前的Output, 倒序来看, 每个Output只有可能传递给之前的Input, 所以不用提前获取全部Inputs, 视频开头弹幕说的是错的
				}
			}
		}
	}
	for _, tx := range txs {
	Workfirst:
		for index, out := range tx.Vouts {
			if out.UnLockScripPubKeyWithAddress(address) { // ← 这个其实就是判断hash是否相同的
				if len(spentTXOutputs) == 0 {
					utxo := &UTXO{tx.TxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash, indexArray := range spentTXOutputs {
						txHashStr := hex.EncodeToString(tx.TxHash)
						if txHashStr == hash {
							var isUnSpentUTXO bool // ← 默认false
							for _, outIndex := range indexArray {
								if index == outIndex {
									isUnSpentUTXO = true
									continue Workfirst
								}
								if isUnSpentUTXO == false {
									utxo := &UTXO{tx.TxHash, index, out}
									unUTXOs = append(unUTXOs, utxo)
								}
							}
						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				}
			}
		}
	}
	blockIterator := blockchain.Iterator() // 返回迭代器
	for {
		block := blockIterator.Next() // ← 如果打印得话, 打印得是最新的区块

		fmt.Println(block)
		fmt.Println("------这部分出现在Blockchain的UnUTXOs函数中------") // ← go语言里面双引号""和单引号''含义不同

		for i := len(block.Txs) - 1; i >= 0; i-- {
			// ↓tx包含以下三个东西
			//  txHash
			tx := block.Txs[i]
			//  Vins
			if tx.IsCoinbaseTransaction() == false { // ← 这行代码说明其是其他正常的transaction
				for _, in := range tx.Vins {
					publicKeyHash := Base58Decode([]byte(address))           // 反编码
					ripemd160Hash := publicKeyHash[1 : len(publicKeyHash)-4] // ← 取公钥Hash部分
					// ↓是否能够解锁
					if in.UnLockRipemd160Hash(ripemd160Hash) {
						key := hex.EncodeToString(in.Txid)
						spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout) // ← 这个"Vout"不是"Vouts", 而是TXInputs里面标注的其继承的交易的索引
						// ↑不用获取全部的Inputs后再核验Output, 因为每个Input只有可能调用自之前的Output, 倒序来看, 每个Output只有可能传递给之前的Input, 所以不用提前获取全部Inputs, 视频开头弹幕说的是错的
					}
				}
			}

			//  Vouts
		work:
			for index, out := range tx.Vouts {
				if out.UnLockScripPubKeyWithAddress(address) {
					if spentTXOutputs != nil {
						// ↓字典也可以range这样遍历: 前面的返回值是每个"键", 后面的返回值是每个"值"
						//for txHash, indexArray := range spentTXOutputs {
						//	//if txHash == hex.EncodeToString(tx.TxHash) {
						//	for _, i := range indexArray {
						//		if index == i && txHash == hex.EncodeToString(tx.TxHash) {
						//			// ↑ 说明这笔钱已经被花费掉了: 已经查找到的Output花费里的某一笔钱, 与保存的该地址的Input记录的Hash相对应
						//			continue
						//		} else {
						//			unUTXOs = append(unUTXOs, out)
						//		}
						//	}
						//}
						if len(spentTXOutputs) != 0 {
							var isSpentUTXO bool // ← 这个修改实在是太牛了, 值得细细品味
							for txHash, indexArray := range spentTXOutputs {
								//if txHash == hex.EncodeToString(tx.TxHash) {
								for _, i := range indexArray {
									if index == i && txHash == hex.EncodeToString(tx.TxHash) {
										// ↑ 说明这笔钱已经被花费掉了: 已经查找到的Output花费里的某一笔钱, 与保存的该地址的Input记录的Hash相对应
										isSpentUTXO = true // ← 当前这笔有没有被消费
										continue work
									}
								}
							}
							if isSpentUTXO == false { // ← 意思是, 如果没有与之对应的Input(就是这笔钱还没花)(现在没花之前肯定也没花), 那就直接记录进账户
								utxo := &UTXO{tx.TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)
							}
						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				}
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		// ↓ 设置退出
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			// ↑上述情况说明hashInt = big.NewInt(0), 即找到对应块了(即其前面没有块了); 故退出(此时理解仅供测试用)
			break
		}
	}
	return unUTXOs
}

// 转账时查找可用的UTXO
func (blockchain *Blockchain) FindSpendableUTXOs(from string, amount int, txs []*Transaction) (int64, map[string][]int) {
	// 1. 先获取所有UTXO
	utxos := blockchain.UnUTXOs(from, txs)
	spendableUTXO := make(map[string][]int)
	// 2. 遍历utxos
	var value int64
	for _, utxo := range utxos {
		value = value + utxo.Output.Value
		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)
		if value >= int64(amount) {
			break
		}
	}
	if value < int64(amount) {
		fmt.Printf("%s's fund is invalid...\n", from)
		os.Exit(1)
	}
	return value, spendableUTXO
}

// ↓挖掘新的区块
func (blockchain *Blockchain) MineNewBlock(from []string, to []string, amount []string) {
	// 1. 建立一笔交易
	//value, _ := strconv.Atoi(amount[0])                           // ← value是转换后的结果, 另一个返回值error是转换过程中的报错
	//tx := NewSimpleTransaction(from[0], to[0], value, blockchain) // ← strconv.Atoi(): 将string格式转为int格式

	// ↓ 下面这三行输出注释掉, 转账消息详情就不会输出了
	//fmt.Println(from)
	//fmt.Println(to)
	//fmt.Println(amount)

	var txs []*Transaction // ← 因为这个txs最终处理后也是个空的txs, 所以产生区块的"Txs:"显示部分为空
	for index, address := range from {
		value, _ := strconv.Atoi(amount[index])                                // ← value是转换后的结果, 另一个返回值error是转换过程中的报错
		tx := NewSimpleTransaction(address, to[index], value, blockchain, txs) // ← strconv.Atoi(): 将string格式转为int格式
		txs = append(txs, tx)
	}

	// 挖矿奖励(奖励一个币)
	tx := TransactionReward(from[0])
	txs = append(txs, tx)

	// 1. 通过相关算法建立交易(Transaction)数组
	var block *Block
	blockchain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockTableName))
		if b != nil {
			// ↓或许hash倒序(最新的在最前, 被获取)排列中的hash
			hash := b.Get([]byte("list"))
			// ↓检索该hash对应的序列化后的结果
			blockBytes := b.Get(hash)
			// ↓反序列化该hash对应的block, 获取具体数据
			block = DeserializeBlock(blockBytes)
		}
		return nil
	})

	// 在建立新的区块前, 对txs进行签名验证

	_txs := []*Transaction{}

	for _, tx := range txs {
		if blockchain.VerifyTransaction(tx, _txs) != true {
			log.Panic("ERROR: Invalid transaction.")
		}
		_txs = append(_txs, tx)
	}

	// 2. 建立新的区块(记得高度+1)
	block = NewBlock(txs, block.Height+1, block.Hash)

	// 3. 存储新区块到数据库
	blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockTableName))
		if b != nil {
			// ↓↓存储数据进区块链
			// ↓更新区块序列化列
			b.Put(block.Hash, block.Serialize()) // b.Put(索引列, 需存储数据) → 在对应的索引列, 存对应的数据
			// ↓更新区块hash排序列
			b.Put([]byte("list"), block.Hash)
			// ↓更新blockchain的Tip
			blockchain.Tip = block.Hash
		}
		return nil
	})
}

// ↓ 查询某个人的总余额
func (blockchain *Blockchain) GetBalance(address string) int64 {
	utxos := blockchain.UnUTXOs(address, []*Transaction{}) // ← 查询余额传空就行
	var amount int64
	for _, utxo := range utxos {
		amount = amount + utxo.Output.Value
	}
	return amount
}

// ↓ 进行数字签名
func (blockchain *Blockchain) SignTransaction(tx *Transaction, private ecdsa.PrivateKey, txs []*Transaction) {
	if tx.IsCoinbaseTransaction() { // ← 如果是创世区块, 啥也不动
		return
	}

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vins {
		prevTX, err := blockchain.FindTransaction(vin.Txid, txs)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.TxHash)] = prevTX
	}

	tx.Sign(private, prevTXs)

}

func (blockchain *Blockchain) FindTransaction(ID []byte, txs []*Transaction) (Transaction, error) {
	for _, tx := range txs {
		if bytes.Compare(tx.TxHash, ID) == 0 {
			return *tx, nil
		}
	}

	bci := blockchain.Iterator()

	for {
		block := bci.Next()
		for _, tx := range block.Txs {
			if bytes.Compare(tx.TxHash, ID) == 0 {
				return *tx, nil
			}
		}
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
	return Transaction{}, nil
}

// 验证数字签名
func (blockchain *Blockchain) VerifyTransaction(tx *Transaction, txs []*Transaction) bool {

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vins {
		prevTX, err := blockchain.FindTransaction(vin.Txid, txs)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.TxHash)] = prevTX
	}

	return tx.Verify(prevTXs)
}

// [string]*TXOutputs
func (blockchain *Blockchain) FindUTXOMap() map[string]*TXOutputs {
	blcIterator := blockchain.Iterator()

	// ↓ 用于存储已经消费(有Inputs)的Hash
	spentableUTXOsMap := make(map[string][]*TXInput)

	utxoMaps := make(map[string]*TXOutputs)

	for {
		block := blcIterator.Next() // ← 拿到最新的区块
		for i := len(block.Txs) - 1; i >= 0; i-- {
			// ↓ 用于存储未花费的TXOutputs
			txOutputs := &TXOutputs{[]*UTXO{}}

			tx := block.Txs[i]

			// 过滤掉没有Inputs的transaction
			if tx.IsCoinbaseTransaction() == false {
				for _, txInput := range tx.Vins {
					txHash := hex.EncodeToString(txInput.Txid)
					spentableUTXOsMap[txHash] = append(spentableUTXOsMap[txHash], txInput)
				}
			}

			txHash := hex.EncodeToString(tx.TxHash)

		NextOutLoop:
			for index, out := range tx.Vouts {
				txInputs := spentableUTXOsMap[txHash]
				if len(txInputs) > 0 {

					isSpent := false

					for _, in := range txInputs {
						outPublicKey := out.Ripemd160Hash
						inPublicKey := in.PublicKey

						// HashPubKey()方法就是视频中的Ripemd160Hash()方法
						if bytes.Compare(outPublicKey, HashPubKey(inPublicKey)) == 0 {
							if index == in.Vout {
								isSpent = true
								continue NextOutLoop
							}
						}
					}
					if isSpent == false {
						utxo := &UTXO{tx.TxHash, index, out}
						txOutputs.UTXOs = append(txOutputs.UTXOs, utxo)
					}
				} else {
					utxo := &UTXO{tx.TxHash, index, out}
					txOutputs.UTXOs = append(txOutputs.UTXOs, utxo)
				}
			}

			// 设置键值对
			utxoMaps[txHash] = txOutputs
		}

		// ↓ 当找到最初的创世区块, 退出
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}

	}

	return utxoMaps
}
