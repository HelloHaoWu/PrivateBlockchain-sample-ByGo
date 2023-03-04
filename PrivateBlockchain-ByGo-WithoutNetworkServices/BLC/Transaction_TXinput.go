package BLC

import "bytes"

// ↓↓记录Txoutput的每笔转账的详情
// ↓这里是表明继承自之前的那一笔交易(Txid)的哪一个TxOutput(Vout)
//
//	加入之前有一笔交易, 交易的Txid为111, 他继承自Txid为000的交易中, 第3个Output"Wang Da Na"(表明接下来要花王大拿的钱), 则其input为
//	Txid: 000, Vout: 3, ScripSig: "Wang Da Na"; 即在0000交易第3个Output的王大拿的数据记录上做修改
//	↑这样做的好处是能实时使用最新的某个用户的状态, 这可太NB了
type TXInput struct {
	//1.引用的状态/继承的状态的交易的ID
	Txid []byte
	//2.在引用的交易中, 所需要引用的钱包状态在引用的交易的Vout的TxOutput里面的索引
	Vout int // 对应交易在Transaction大记录的Txoutput中的索引为几

	//3.用户名(已弃用)
	//ScriptSig string //用户签名(将要消费谁的钱)

	//3.数字签名(正式版使用)
	Signature []byte

	//4.公钥(正式版使用) ← 钱包里面的公钥
	PublicKey []byte
}

// ↑由于某个账户的"钱包", 可能分布在多个Transaction中, 所以可能涉及多调用的问题; 这样省去了input的问题, 而是直接调用对应transaction的output来计算input

// ↓传一个地址过来, 判断它的txInput的名字(ScriptSig)是否和要查询的名字对应
func (txInput *TXInput) UnLockRipemd160Hash(ripemd160Hash []byte) bool {
	publicKey := HashPubKey(txInput.PublicKey)          // hash转160加密
	return bytes.Compare(publicKey, ripemd160Hash) == 0 // 说明刚好对应
}
